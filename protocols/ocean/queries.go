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

const AQUARIUS_URL_DDO = "https://aquarius.oceanprotocol.com/api/v1/aquarius/assets/ddo/"
const PURGATORY_ASSETS = "https://raw.githubusercontent.com/oceanprotocol/list-purgatory/main/list-assets.json"
const PURGATORY_ACCOUNTS = "https://raw.githubusercontent.com/oceanprotocol/list-purgatory/main/list-assets.json"

// graphQuery gets most blockchain data from Ocean's GraphQL instance.
func graphQuery(query string, respContainer interface{}, debug bool) (err error) {
	// create a client (safe to share across requests)
	client := graphql.NewClient(OCEAN_SUBGRAPH_MAINNET)
	if debug {
		client.Log = func(s string) { log.Println(s) }
	}
	req := graphql.NewRequest(query)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	if err = client.Run(ctx, req, respContainer); err != nil {
		log.Fatal(err)
		return
	}
	return nil
}

// aquariusError is needed to tell the upper layer more nuanced errors, like
// whether it was 404 not found or 503 service unavailable
type aquariusError struct {
	StatusCode int
	Body       []byte
}

func (ae *aquariusError) Error() string {
	return fmt.Sprintf("Aquarius error: %v %v", ae.StatusCode, ae.Body)
}

// aquariusQuery gets additional data like purgatory status and a dataset
// description from Aquarius. The argument datatokenAddress must be the 0x...
// Ethereum address, which will be stripped of its 0x prefix to produce the DID.
func aquariusQuery(datatokenAddress string) (ddo *DecentralizedDataObject, err error) {
	did := strings.TrimLeft(datatokenAddress, "0x")
	resp, err := http.Get(fmt.Sprintf("%sdid:op:%s", AQUARIUS_URL_DDO, did)) // https://multiaqua.oceanprotocol.com/api/v1/aquarius/assets/ddo/did:op:0f5A4C51Dd71C7FB8D5D61e5B56C996681e4302F
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Ocean Aquarius did not return a OK to our request for did %s. Status Code: %s, Body: %s", did, resp.Status, body)
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

// purgAccounts gets a list of assets in purgatory from a fixed URL on Github
// and parses the JSON. Then it transforms the list into a map[string] for easy
// lookup
func purgAccounts() (purgatoryMap map[string]string, err error) {
	var pa []PurgatoryAccount
	purgatoryMap = make(map[string]string)
	resp, err := http.Get(PURGATORY_ACCOUNTS)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &pa)
	if err != nil {
		return
	}

	for _, a := range pa {
		purgatoryMap[a.Address] = a.Reason
	}
	return purgatoryMap, nil
}
