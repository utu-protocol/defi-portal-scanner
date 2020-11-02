package uniswap

import (
	"context"
	"fmt"
	"log"

	"github.com/machinebox/graphql"
)

/**

- Positive BALANCE

liquidityTokenBalance

## AMM (3 most popular)
- liquidity
    - (liquidityPool)
    - (accounts) - trades - (liquidityPool)  --> tradeVolume
        - (liquidityProvider) - add/remove - (liquidity) --> size


## Borrow/lending (important for the business story)
- borrowing/lending

*/

// Start start processing uniswap data
func Start(theGraphEndpoint string) {
	// create a client (safe to share across requests)
	client := graphql.NewClient("https://machinebox.io/graphql")

	// We want the transactions
	req := graphql.NewRequest(`{
		transactions(first: 100){
			id
			blockNumber
			timestamp
			swaps {
				id
				pair {
					id
					token0 {
						name
					}
					token1 {
						name
					}
				}
				sender
				to
				amountUSD
			}
			mints {
				id
				liquidity
				feeTo
				to
			}
			burns {
				id
			}
		}
	}`)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var respData []map[string]interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}
	fmt.Println(respData)

}
