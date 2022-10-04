package ocean

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"crypto/sha256"

	"github.com/machinebox/graphql"
	"github.com/utu-crowdsale/defi-portal-scanner/utils"
)

const OCEAN_ERC20_ADDRESS = "0x967da4048cd07ab37855c090aaf366e4ce1b9f48"
const OCEAN_SUBGRAPH_MAINNET = "https://v4.subgraph.mainnet.oceanprotocol.com/subgraphs/name/oceanprotocol/ocean-subgraph"
const CHAIN_ID = 1

const AQUARIUS_URL_DDO = "https://v4.aquarius.oceanprotocol.com/api/aquarius/assets/ddo/"
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

	// run it and capture the response
	if err = client.Run(context.Background(), req, respContainer); err != nil {
		log.Fatal(err)
		return
	}
	return nil
}

// aquariusError is needed to tell the upper layer more nuanced errors, like
// whether it was 404 not found or 503 service unavailable
type aquariusError struct {
	RequestedDID string
	StatusCode   int
	Body         []byte
}

func (ae *aquariusError) Error() string {
	return fmt.Sprintf("Aquarius error while requesting did %s: %d %s", ae.RequestedDID, ae.StatusCode, ae.Body)
}

// aquariusQuery gets additional data like purgatory status and a dataset
// description from Aquarius. The argument datatokenAddress must be the 0x...
// Ethereum address, which will be stripped of its 0x prefix to produce the DID.
// IMPORTANT: The DID is made from a checksummed Ethereum address, not lowercase!
func aquariusQuery(erc721Address string) (ddo *DecentralizedDataObject, err error) {
	erc721AddressChecksumed := utils.ChecksumAddress(erc721Address)
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%s%v", erc721AddressChecksumed, CHAIN_ID)))
	did := hex.EncodeToString(h.Sum(nil))
	requestURL := fmt.Sprintf("%sdid:op:%s", AQUARIUS_URL_DDO, did)
	log.Println(requestURL)
	resp, err := http.Get(requestURL)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		return nil, &aquariusError{
			RequestedDID: did,
			StatusCode:   resp.StatusCode,
			Body:         body,
		}
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
