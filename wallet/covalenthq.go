package wallet

import (
	"encoding/json"
	"math/big"

	log "github.com/sirupsen/logrus"
)

type CovalentHQResponse struct {
	Data struct {
		ChainID int                      `json:"chain_id"`
		Items   []CovalentHQResponseItem `json:"items"`
	} `json:"data"`
}

type CovalentHQResponseItem struct {
	Decimals int    `json:"contract_decimals"`
	Name     string `json:"contract_name"`
	Symbol   string `json:"contract_ticker_symbol"`
	Address  string `json:"contract_address"`
	Balance  string `json:"balance"`
}

func scanCovalentBalance(apiURL string) (*Response, error) {
	responseBytes, err := httpGet(apiURL)
	if err != nil {
		log.Errorf("Cannot fetch token balances err=%s", err.Error())
		return nil, err
	}
	var response CovalentHQResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		log.Errorf("Cannot unmarshall response to a struct err=%s", err.Error())
		return nil, err
	}
	return &Response{
		APIResponse: &response,
	}, nil
}

func mapCovalentResponse(r *Response, network string) (balances []Balance) {
	response := r.APIResponse.(*CovalentHQResponse)
	for _, item := range response.Data.Items {
		tokenBalance, ok := new(big.Int).SetString(item.Balance, 10)
		if !ok {
			log.Errorf("cannot parse balance")
			continue
		}
		balance := Balance{
			Address: item.Address,
			Balance: formatBalance(tokenBalance, item.Decimals),
			Name:    item.Name,
			Symbol:  item.Symbol,
			Network: network,
		}
		balances = append(balances, balance)
	}
	return
}
