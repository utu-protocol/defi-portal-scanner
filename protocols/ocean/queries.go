package ocean

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/machinebox/graphql"
)

const OCEAN_ERC20_ADDRESS = "0x967da4048cd07ab37855c090aaf366e4ce1b9f48"
const OCEAN_SUBGRAPH_MAINNET = "https://subgraph.mainnet.oceanprotocol.com/subgraphs/name/oceanprotocol/ocean-subgraph"

const AQUARIUS_URL_DDO = "https://multiaqua.oceanprotocol.com/api/v1/aquarius/assets/ddo/"
const PURGATORY_ASSETS = "https://github.com/oceanprotocol/list-purgatory/blob/main/list-assets.json"
const PURGATORY_ACCOUNTS = "https://github.com/oceanprotocol/list-purgatory/blob/main/list-assets.json"

// graphQuery gets most blockchain data from Ocean's GraphQL instance.
func graphQuery(query string) (respData map[string]interface{}, err error) {
	// create a client (safe to share across requests)
	client := graphql.NewClient(OCEAN_SUBGRAPH_MAINNET)

	req := graphql.NewRequest(query)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	if err = client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
		return
	}
	return respData, nil
}

// aquariusQuery gets additional data like purgatory status and a dataset
// description from Aquarius. The argument datatokenAddress must be the 0x...
// Ethereum address, which will be stripped of its 0x prefix to produce the DID.
func aquariusQuery(datatokenAddress string) (ddo *DecentralizedDataObject, err error) {
	did := strings.TrimLeft(datatokenAddress, "0x")
	resp, err := http.Get(fmt.Sprintf("%s%s", AQUARIUS_URL_DDO, did))
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	ddo = new(DecentralizedDataObject)
	err = json.Unmarshal(body, ddo)
	return
}

type PurgatoryAsset struct {
	DID    string
	Reason string
}
type PurgatoryAccount struct {
	Address string
	Reason  string
}

// purgAssets gets a list of assets in purgatory from a fixed URL on Github and parses the JSON.
func purgAssets() (pa []PurgatoryAsset, err error) {
	resp, err := http.Get(PURGATORY_ASSETS)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &pa)
	return
}

// purgAccounts gets a list of assets in purgatory from a fixed URL on Github and parses the JSON.
func purgAccounts() (pa []PurgatoryAccount, err error) {
	resp, err := http.Get(PURGATORY_ACCOUNTS)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &pa)
	return
}
