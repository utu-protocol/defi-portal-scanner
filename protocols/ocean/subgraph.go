package ocean

import (
	"context"
	"fmt"
	"log"

	"github.com/machinebox/graphql"
)

const OCEAN_ERC20_ADDRESS = "0x967da4048cd07ab37855c090aaf366e4ce1b9f48"
const OCEAN_SUBGRAPH_MAINNET = "https://subgraph.mainnet.oceanprotocol.com/subgraphs/name/oceanprotocol/ocean-subgraph"

// pools
func pools() {
	// create a client (safe to share across requests)
	client := graphql.NewClient(OCEAN_SUBGRAPH_MAINNET)

	// We want the transactions
	req := graphql.NewRequest(`{
		pools{
			controller,
			totalSwapVolume,
			datatokenReserve,
			datatokenAddress,
			spotPrice,
			consumePrice,
			tokens{denormWeight, tokenAddress, balance}
		  }
	}`)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var respData map[string]interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}
	fmt.Println(req, respData)

}
