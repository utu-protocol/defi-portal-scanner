package ocean

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/barkimedes/go-deepcopy"
)

func paginatedGraphQuery(baseQuery string, respContainer pageEmptiable) (pages []interface{}, err error) {
	i := 0
	for {
		var query string
		if i == 0 {
			query = fmt.Sprintf(baseQuery, ",first:1000")
		} else {
			s := fmt.Sprintf(",first:1000,skip:%d", i)
			query = fmt.Sprintf(baseQuery, s)
		}

		err = graphQuery(query, respContainer)
		if err != nil {
			log.Println("Error while querying GraphQL for datatokens", err)
			return
		}
		if respContainer.IsEmpty() {
			break
		}

		page := deepcopy.MustAnything(respContainer)
		pages = append(pages, page)
		i += 1000
	}
	return
}

// pipeline makes queries to Aquarius, the github repos for things in purgatory,
// and builds up an internal state.
func pipeline(log *log.Logger) (err error) {
	// Get basic data about Datatokens, and how many times they were consumed.
	baseDatatokensQuery := `{datatokens(orderBy:name%s) {
		symbol
		name
		address
		orders{
		  consumer{
			id
		  }
		}
		orderCount
	  }}`

	pages, err := paginatedGraphQuery(baseDatatokensQuery, new(DatatokenResponsePage))
	if err != nil {
		log.Println("Error connecting to GraphQL", err)
		return err
	}
	/*
		"datatokens": [
			{
				"address": "0x028e0b27a39ff92fd30b4b8c310ea745f309ccf3",
				"name": "Brave Nautilus Token",
				"orderCount": "3",
				"orders": [
					{
					"consumer": {
						"id": "0x1bb7951ba30eda67bf3e5d851fe5e0e6a01a14b5"
					}
					},
					{
					"consumer": {
						"id": "0x4ba10551d7b76b30369e9ef8d27966e19dcc786b"
					}
					},
					{
					"consumer": {
						"id": "0xb40156f51103ebaa842590ce51dd2cd0a9e83cda"
					}
					}
				],
				"symbol": "BRANAU-77"
			}
			]
	*/
	var dtr []DatatokenResponse
	for _, p := range pages {
		page := p.(*DatatokenResponsePage)
		dtr = append(dtr, page.Flatten()...)
	}

	var datatokens []*Datatoken
	for _, datatokenResponse := range dtr {
		dt, err := NewDataTokenFromDatatokenResponse(datatokenResponse)
		if err != nil {
			log.Println("Could not create a DataToken internal class from a datatokenResponse", err)
			return err
		}
		datatokens = append(datatokens, dt)
	}

	// Now pull up info on all Pools. Since we receive Pools as a list, which is
	// not so easy to search through, we transform it into a map of
	// datatokenAddress -> PoolGraphQLResponse
	basePoolsQuery := `{pools (where: {datatokenAddress_not: ""}, orderBy: oceanReserve, orderDirection:desc%s ) {
		id
		controller
		totalSwapVolume
		datatokenAddress
		datatokenReserve
		oceanReserve
	}}`

	respPool := new(PoolResponsePage)
	pages, err = paginatedGraphQuery(basePoolsQuery, respPool)
	if err != nil {
		log.Println("Error while querying GraphQL for pools", err)
		return
	}
	var pr []PoolGraphQLResponse
	for _, p := range pages {
		page := p.(*PoolResponsePage)
		pr = append(pr, page.Flatten()...)
	}
	pm := make(map[string]*Pool)
	for _, pGrQlResp := range pr {
		pool, err := pGrQlResp.toPool()
		if err != nil {
			log.Println("Error while transforming PoolGraphQLResponse to Pool struct", err)
			return err
		}

		pm[checksumAddress(pGrQlResp.DatatokenAddress)] = pool
	}

	// We have Pools and Datatokens, so we can now construct Assets. A Datatoken
	// may not have a Pool, or even a DDO associated with it. But we're only
	// interested in Datatokens with Pools and DDOs.
	var assets []*Asset
	for _, dt := range datatokens {
		// If Datatoken does not have corresponding Pool, immediately skip it (don't bother Aquarius)
		pool := pm[dt.Address]
		if pool == nil {
			log.Printf("%s does not have a corresponding Pool, skipping", dt.Address)
			continue
		}

		// In practice, Aquarius only knows about Datatokens which have Pools.
		var purgatoryStatus bool
		ddo, err := aquariusQuery(dt.Address)
		if err != nil {
			log.Printf("%s is not known by Aquarius, skipping: %s", dt.Address, err)
			continue
		}

		purgatoryStatus, err = strconv.ParseBool(ddo.IsInPurgatory)
		if err != nil {
			return err
		}
		asset := &Asset{
			Pool:        pool,
			Datatoken:   dt,
			PublishedBy: pool.Controller,
			Purgatory:   purgatoryStatus,
			Consumed:    dt.OrderCount,
		}
		assets = append(assets, asset)
	}

	j, err := json.MarshalIndent(assets, "", "\t")
	if err != nil {
		return
	}
	fmt.Println(string(j))
	fmt.Println("len(assets)", len(assets))
	return
}

type pageEmptiable interface {
	IsEmpty() bool
}

type DatatokenResponsePage struct {
	Datatokens []DatatokenResponse
}

func (dt *DatatokenResponsePage) IsEmpty() bool {
	return len(dt.Datatokens) == 0
}

func (dt *DatatokenResponsePage) Flatten() []DatatokenResponse {
	return dt.Datatokens
}

type PoolResponsePage struct {
	Pools []PoolGraphQLResponse
}

func (pr *PoolResponsePage) IsEmpty() bool {
	return len(pr.Pools) == 0
}

func (pr *PoolResponsePage) Flatten() []PoolGraphQLResponse {
	return pr.Pools
}
