package collector

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/utu-crowdsale/defi-portal-scanner/config"
	"github.com/utu-crowdsale/defi-portal-scanner/wallet"
	"gopkg.in/robfig/cron.v2"
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
	go balanceReqProcessor()
	go walletProcessor(cfg.UTUTrustAPI)
	c := cron.New()
	c.AddFunc("@every 24h", func() {
		go wallet.ScanCached(walletsChan)
	})
	c.Start()
}

func ScanTokensBalances(address string) {
	log.Infof("Received a request to scan tokens balances for %s", address)
	balancesRequestQueue <- &BalanceRequest{
		Address: strings.ToLower(address),
	}
}

func balanceReqProcessor() {
	for {
		req, more := <-balancesRequestQueue
		log.Infof("received request to scan token balances for %v", req.Address)
		if !more {
			log.Info("no more requests to scan balances for")
			break
		}
		wallet.Scan(req.Address, walletsChan)
	}
}

func walletProcessor(cfg config.TrustEngineSchema) {
	utuCli := NewUTUClient(cfg)
	for {
		wallet, more := <-walletsChan
		log.Infof("received wallet data for %v", wallet.Address)
		if !more {
			log.Info("no more wallet data")
			break
		}
		for _, balance := range wallet.Balances {
			trustRelationship := toTrustRelationship(wallet, &balance)
			if err := utuCli.PostRelationship(trustRelationship); err != nil {
				log.Error("error posting relationship:", err)
			}
		}
	}
}

func toTrustRelationship(wallet *wallet.Wallet, balance *wallet.Balance) *TrustRelationship {
	r := NewTrustRelationship()
	r.Type = "has"
	r.SourceCriteria = NewTrustEntity(wallet.Address)
	r.SourceCriteria.Ids["address"] = wallet.Address
	r.SourceCriteria.Type = "address"
	r.TargetCriteria = NewTrustEntity(balance.Symbol)
	r.TargetCriteria.Ids["address"] = balance.Address
	r.TargetCriteria.Type = "erc20"
	r.TargetCriteria.Name = balance.Symbol
	r.Properties["balance"] = balance.Balance
	r.Properties["network"] = balance.Network
	return r
}
