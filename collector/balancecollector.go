package collector

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/utu-crowdsale/defi-portal-scanner/config"
	"github.com/utu-crowdsale/defi-portal-scanner/wallet"
)

type BalanceRequest struct {
	Address string
	Tokens  []string
}

var (
	balancesRequestQueue chan *BalanceRequest
	walletsChan          chan *wallet.Wallet
)

func init() {
	balancesRequestQueue = make(chan *BalanceRequest)
	walletsChan = make(chan *wallet.Wallet)
}

func BalanceCollectorReady(cfg config.Schema) {
	go balanceReqProcessor(cfg)
	go walletProcessor(cfg)
}

func ScanTokensBalances(cfg config.Schema, address string, tokens []string) {
	balancesRequestQueue <- &BalanceRequest{
		Address: strings.ToLower(address),
		Tokens:  tokens,
	}
}

func balanceReqProcessor(cfg config.Schema) {
	for {
		req, more := <-balancesRequestQueue
		log.Infof("received request to scan token balances for %v", req.Address)
		if !more {
			log.Info("no more requests to scan balances for")
			break
		}
		wallet.Scan(cfg, req.Address, req.Tokens, walletsChan)
	}
}

func walletProcessor(cfg config.Schema) {
	for {
		wallet, more := <-walletsChan
		log.Infof("received wallet data for %v", wallet.Address)
		if !more {
			log.Info("no more wallet data")
			break
		}
		fmt.Println(*wallet)
	}
}
