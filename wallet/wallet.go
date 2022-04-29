package wallet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"github.com/utu-crowdsale/defi-portal-scanner/config"
)

type tokenData struct {
	Name      string    `json:"name"`
	LogoURI   string    `json:"logoURI"`
	Keywords  []string  `json:"keywords"`
	Timestamp time.Time `json:"timestamp"`
	Tokens    []token   `json:"tokens"`
}

type token struct {
	ChainID  int    `json:"chainId"`
	Address  string `json:"address"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
	LogoURI  string `json:"logoURI"`
}

func (t *token) addressToLowercase() {
	t.Address = strings.ToLower(t.Address)
}

type Wallet struct {
	Address  string    `json:"address"`
	Balances []Balance `json:"balances"`
}

type Balance struct {
	Address  string `json:"address"`
	Balance  string `json:"balance"`
	Decimals int    `json:"decimals"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
}

type balanceRequest struct {
	address string
	tokens  interface{}
}

var (
	tokenDataByAddress   map[string]token
	contractAddresses    []string
	balancesRequestQueue chan *balanceRequest
)

func init() {
	balancesRequestQueue = make(chan *balanceRequest)
}

func Ready(cfg config.Schema) {
	go scan(cfg)
}

func ScanTokensBalances(cfg config.Schema, address string, tokens []string) {
	loadTokensData(cfg)
	var addresses interface{}
	if len(tokens) > 0 {
		addressesLowerCase := make([]string, 0)
		for _, addr := range tokens {
			addressesLowerCase = append(addressesLowerCase, strings.ToLower(addr))
		}
		addresses = addressesLowerCase
	} else if len(contractAddresses) > 0 {
		addresses = contractAddresses
	} else {
		addresses = "DEFAULT_TOKENS"
	}
	balancesRequestQueue <- &balanceRequest{
		address: strings.ToLower(address),
		tokens:  addresses,
	}
}

func scan(cfg config.Schema) {
	for {
		req, more := <-balancesRequestQueue
		log.Infof("received request to scan token balances for %v", req.address)
		if !more {
			log.Info("no more requests to scan balances for")
			break
		}
		alchemyResponse, err := getBalances(cfg, req.address, req.tokens)
		if err != nil {
			log.Error("Couldn't fetch token balances for %s", req.address, err)
			return
		}
		fmt.Println(*alchemyResponse)
		balances := make([]Balance, 0)
		for _, token := range alchemyResponse.Result.TokenBalances {
			balanceHex := common.FromHex(token.TokenBalance)
			tokenBalance := new(big.Int).SetBytes(balanceHex)
			if token.Error != nil || tokenBalance.BitLen() == 0 {
				continue
			}
			tokenData := tokenDataByAddress[token.ContractAddress]
			balance := Balance{
				Address:  tokenData.Address,
				Balance:  tokenBalance.String(),
				Decimals: tokenData.Decimals,
				Name:     tokenData.Name,
				Symbol:   tokenData.Symbol,
			}
			balances = append(balances, balance)
		}
		if len(balances) > 0 {
			wallet := &Wallet{
				Address:  req.address,
				Balances: balances,
			}
			fmt.Println(wallet)
		}
	}
}

func loadTokensData(cfg config.Schema) {
	if len(tokenDataByAddress) != 0 {
		return
	}
	log.Infof("Reading tokens metadata from %s", cfg.TokensDataFile)
	tokenDataByAddress = make(map[string]token)
	erc20tokensJson, err := ioutil.ReadFile(cfg.TokensDataFile)
	if err != nil {
		log.Errorf("Cannot read tokens metadata file to load default tokens list, %s", err.Error())
		return
	}
	var tokenData tokenData
	err = json.Unmarshal(erc20tokensJson, &tokenData)
	if err != nil {
		log.Errorf("Cannot unmarshal tokens metadata, %s", err.Error())
		return
	}
	for _, token := range tokenData.Tokens {
		token.addressToLowercase()
		tokenDataByAddress[token.Address] = token
	}
	contractAddresses = make([]string, 0, len(tokenDataByAddress))
	for key := range tokenDataByAddress {
		contractAddresses = append(contractAddresses, key)
	}
}
