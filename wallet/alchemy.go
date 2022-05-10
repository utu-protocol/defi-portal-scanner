package wallet

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/utu-crowdsale/defi-portal-scanner/config"
)

const contentType = "application/json"

type alchemyRequest struct {
	JsonRpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type alchemyResponse struct {
	Id      string `json:"id"`
	JsonRpc string `json:"jsonrpc"`
	Result  struct {
		Address       string         `json:"address"`
		TokenBalances []tokenBalance `json:"tokenBalances"`
	} `json:"result"`
}

type tokenBalance struct {
	ContractAddress string      `json:"contractAddress"`
	TokenBalance    string      `json:"tokenBalance"`
	Error           interface{} `json:"error"`
}

func getBalances(cfg config.Schema, address string, addresses interface{}) (*alchemyResponse, error) {
	response, err := doPost(cfg, "alchemy_getTokenBalances", []interface{}{address, addresses})
	if err != nil {
		return nil, err
	}
	var aResponse alchemyResponse
	err = json.Unmarshal(response, &aResponse)
	if err != nil {
		return nil, err
	}
	return &aResponse, nil
}

func doPost(cfg config.Schema, method string, params []interface{}) ([]byte, error) {
	requestStructure := alchemyRequest{
		JsonRpc: "2.0",
		Method:  method,
		Params:  params,
	}
	payload, err := bodyToPayload(requestStructure)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(cfg.AlchemyAPI.URL, contentType, payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func bodyToPayload(request alchemyRequest) (*bytes.Buffer, error) {
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewBuffer([]byte(requestBytes))
	return payload, nil
}
