package collector

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	"github.com/utu-crowdsale/defi-portal-scanner/config"
)

func OceanPoolEventSubscribe(ethConfig config.EthereumSchema, oceanPools []config.Address) {
	log.Info("Starting collector for Ocean Pools")

	// First, cast our Address into a ethereum library common.Address
	var oceanPoolsFilterList []common.Address
	for a := range oceanPools {
		oceanPoolsFilterList = append(oceanPoolsFilterList, common.HexToAddress(string(oceanPools[a])))
	}

	// Connect to the Ethereum node
	client, err := ethclient.Dial(ethConfig.WssURL)
	if err != nil {
		return
	}

	// Listen only to events related to the listed Ocean Pools
	log.Infof("Subscribing to events from %v Ocean Pools", len(oceanPoolsFilterList))
	query := ethereum.FilterQuery{Addresses: oceanPoolsFilterList}
	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			changeset, _ := ParseLog(&vLog, client)
			log.Info(changeset)
		}
	}
}
