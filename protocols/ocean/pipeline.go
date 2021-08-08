package ocean

import (
	"log"
	"strconv"
)

// pipeline makes queries to Aquarius, the github repos for things in purgatory,
// and builds up an internal state.
func pipeline(log *log.Logger) (err error) {
	// Get basic data about Datatokens, and how many times they were consumed.
	query := `{datatokens {
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
	resp := new(datatokenResponse)
	err = graphQuery(query, resp)
	if err != nil {
		log.Println("Error while querying GraphQL for datatokens", err)
		return
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
	dt := make([]Datatoken, len(resp.Datatokens))
	for k, v := range resp.Datatokens {
		orderCount, err := strconv.ParseUint(v.OrderCount, 10, 64)
		if err != nil {
			return err
		}
		dt[k] = Datatoken{
			Address:    v.Address,
			Name:       v.Name,
			Symbol:     v.Symbol,
			OrderCount: uint(orderCount),
		}
	}

	// fmt.Printf("%+v\n", dt)

	// Now pull up info on all Pools
	query = `{pools (where: {datatokenAddress_not: ""}, orderBy: oceanReserve, orderDirection:desc ) {
		id
		controller
		totalSwapVolume
		datatokenAddress
		datatokenReserve
		oceanReserve
	}}`
	type poolResponse struct {
		Pools []struct {
			Controller       string
			DatatokenAddress string
			DatatokenReserve string
			ID               string
			OceanReserve     string
			TotalSwapVolume  string
		}
	}
	respPool := new(poolResponse)
	err = graphQuery(query, respPool)
	if err != nil {
		log.Println("Error while querying GraphQL for pools", err)
		return
	}

	return
}
