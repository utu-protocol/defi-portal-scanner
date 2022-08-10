package wallet

import (
	"encoding/json"
	"math/big"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type BlockscoutResponse struct {
	Result []BlocksoutResponseItem `json:"result"`
}

type BlocksoutResponseItem struct {
	Decimals string `json:"decimals"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Address  string `json:"contractAddress"`
	Balance  string `json:"balance"`
	Type     string `json:"type"`
}

func scanBlockscoutBalance(apiURL string) (*Response, error) {
	responseBytes, err := httpGet(apiURL)
	if err != nil {
		log.Errorf("Cannot fetch token balances err=%s", err.Error())
		return nil, err
	}
	var response BlockscoutResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		log.Errorf("Cannot unmarshall response to a struct err=%s", err.Error())
		return nil, err
	}
	return &Response{
		APIResponse: &response,
	}, nil
}

func mapToBalance(r *Response, network string) (balances []Balance) {
	response := r.APIResponse.(*BlockscoutResponse)
	for _, item := range response.Result {
		if item.Type != "ERC-20" {
			continue
		}
		tokenBalance, ok := new(big.Int).SetString(item.Balance, 10)
		if !ok {
			log.Errorf("cannot parse balance")
			continue
		}
		decimalsNumber, err := strconv.Atoi(item.Decimals)
		if err != nil {
			log.Errorf("cannot parse decimals")
			continue
		}
		balance := Balance{
			Address: item.Address,
			Balance: formatBalance(tokenBalance, decimalsNumber),
			Name:    item.Name,
			Symbol:  item.Symbol,
			Network: network,
		}
		balances = append(balances, balance)
	}
	return
}
