package ocean

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/barkimedes/go-deepcopy"
)

func paginatedGraphQuery(baseQuery string, respContainer pageEmptiable) (pages []interface{}, err error) {
	pageSize := 100
	i := 0
	for {
		var query string
		if i == 0 {
			pagingConstraint := fmt.Sprintf(",first:%d", pageSize)
			query = fmt.Sprintf(baseQuery, pagingConstraint)
		} else {
			pagingConstraint := fmt.Sprintf(",first:%d,skip:%d", pageSize, i)
			query = fmt.Sprintf(baseQuery, pagingConstraint)
		}

		err = graphQuery(query, respContainer, false)
		if err != nil {
			log.Println("Error while querying GraphQL for datatokens", err)
			return
		}
		if respContainer.IsEmpty() {
			break
		}

		page := deepcopy.MustAnything(respContainer)
		pages = append(pages, page)
		i += pageSize
	}
	return
}

func PipelineAll(log *log.Logger) (*PipelineResult, error) {

	// Get basic data about Datatokens, and how many times they were consumed.
	baseDatatokensQuery := `{datatokens: tokens(where: {isDatatoken: true} orderBy:name%s) {
			symbol
			name
			address
			nft {
				  creator
				address
			}
			orders{
			  tx
			  amount
			  timestamp: createdTimestamp
			  user: consumer{
				id
			  }
			  price: lastPriceValue
			}
			orderCount
		  }}`

	pages, err := paginatedGraphQuery(baseDatatokensQuery, new(DatatokenResponsePage))
	if err != nil {
		log.Println("Error connecting to GraphQL", err)
		return nil, err
	}

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
			return nil, err
		}
		datatokens = append(datatokens, dt)
	}

	// We have Pools and Datatokens, so we can now construct Assets. A Datatoken
	// may not have a Pool, or even a DDO associated with it. But we're only
	// interested in Datatokens with Pools and DDOs.
	assets := make([]*Asset, 0)
	for _, dt := range datatokens {

		// In practice, Aquarius only knows about Datatokens which have Pools.
		ddo, err := aquariusQuery(dt.NFT.NFTAddress)
		if err != nil {
			log.Printf("%s is not known by Aquarius, skipping: %s", dt.Address, err)
			continue
		}

		asset := &Asset{
			Name:               ddo.Metadata.Name,
			Description:        ddo.Metadata.Description,
			DID:                ddo.ID,
			Datatoken:          dt,
			PublishedBy:        ddo.Metadata.Description,
			PublishedByAddress: dt.NFT.Creator,
			Purgatory:          ddo.Purgatory.State,
			Consumed:           dt.OrderCount,
			Tags:               ddo.Metadata.Tags,
			Categories:         ddo.Metadata.Categories,
		}
		assets = append(assets, asset)
	}

	addresses, err := pipelineUsers(dtr, log)
	if err != nil {
		return nil, err
	}

	_, err = json.MarshalIndent(assets, "", "\t")
	if err != nil {
		return nil, err
	}
	return &PipelineResult{
		Assets:    assets,
		Addresses: addresses,
	}, nil

}

type PipelineResult struct {
	Assets    []*Asset
	Addresses []*Address
}

// PipelineAll makes queries to Aquarius, the github repos for things in purgatory,
// and builds up an internal state.
func pipelineUsers(datatokensResponse []DatatokenResponse, log *log.Logger) (addresses []*Address, err error) {
	// First, get a list of Accounts in Purgatory from Github.
	purgatoryMap, err := purgAccounts()
	if err != nil {
		log.Println("Error getting list of accounts in purgatory", err)
		return
	}

	ordersByUser := make(map[string][]OrderWrapper)
	for _, dd := range datatokensResponse {
		wrapper := OrderWrapper{
			Token: OrderToken{
				ID:     dd.Address,
				Symbol: dd.Symbol,
				Name:   dd.Name,
			}}
		for _, oo := range dd.Orders {
			wrapper.TxHash = oo.Tx
			wrapper.Amount = oo.Amount
			wrapper.Timestamp = oo.Timestamp
			wrapper.Price = oo.Price
			ordersByUser[oo.User.ID] = append(ordersByUser[oo.User.ID], wrapper)
		}
	}

	for user, orders := range ordersByUser {
		u, err := NewAddressFromUserResponse(user, orders, purgatoryMap)
		if err != nil {
			log.Println("Could not create a Address internal class from a userResponse", err)
			return nil, err
		}
		addresses = append(addresses, u)
	}
	return addresses, nil
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
