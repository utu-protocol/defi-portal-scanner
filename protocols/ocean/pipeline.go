package ocean

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

// pipeline makes queries to Aquarius, the github repos for things in purgatory,
// and builds up an internal state.
func pipeline(log *log.Logger) (err error) {
	// Get basic data about Datatokens, and how many times they were consumed.
	baseDatatokensQuery := `{datatokens(orderBy:name,first:1000%v) {
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
	type datatokenResponse struct {
		Datatokens []struct {
			Address    string
			Name       string
			OrderCount string
			Orders     []struct {
				Consumer struct {
					ID string
				}
			}
			Symbol string
		}
	}
	var datatokens []*Datatoken
	i := 0
	for {
		var query string
		resp := new(datatokenResponse)
		if i == 0 {
			query = fmt.Sprintf(baseDatatokensQuery, "")
		} else {
			s := fmt.Sprintf(",skip:%v", i)
			query = fmt.Sprintf(baseDatatokensQuery, s)
		}
		err = graphQuery(query, resp)
		if err != nil {
			log.Println("Error while querying GraphQL for datatokens", err)
			return
		}
		if len(resp.Datatokens) == 0 {
			break
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
		for _, v := range resp.Datatokens {
			orderCount, err := strconv.ParseUint(v.OrderCount, 10, 64)
			if err != nil {
				return err
			}
			dt, err := NewDataToken(v.Address, v.Name, v.Symbol, uint64(orderCount))
			if err != nil {
				return err
			}
			datatokens = append(datatokens, dt)
		}
		i += 1000
	}

	// Now pull up info on all Pools. Since we receive Pools as a list, which is
	// not so easy to search through, we transform it into a map of
	// datatokenAddress -> PoolGraphQLResponse
	basePoolsQuery := `{pools (where: {datatokenAddress_not: ""}, orderBy: oceanReserve, orderDirection:desc ) {
		id
		controller
		totalSwapVolume
		datatokenAddress
		datatokenReserve
		oceanReserve
	}}`
	type poolResponse struct {
		Pools []PoolGraphQLResponse
	}
	respPool := new(poolResponse)
	err = graphQuery(basePoolsQuery, respPool)
	if err != nil {
		log.Println("Error while querying GraphQL for pools", err)
		return
	}
	pm := make(map[string]*Pool)
	for _, pGrQlResp := range respPool.Pools {
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

	fmt.Println("len(assets)", len(assets))
	j, err := json.MarshalIndent(assets, "", "\t")
	if err != nil {
		return
	}
	fmt.Println(string(j))
	return
}
