/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/utu-crowdsale/defi-portal-scanner/collector"
	"github.com/utu-crowdsale/defi-portal-scanner/utils"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan ADDRESS",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(1),
	Long:  ``,
	Run:   scan,
}

func init() {
	rootCmd.AddCommand(scanCmd)
}

func scan(cmd *cobra.Command, args []string) {

	log.SetFormatter(&utils.EmojiLogFormatter{})
	if debug {
		// Only log the warning severity or above.
		log.SetLevel(log.DebugLevel)
	}

	client, err := ethclient.Dial(settings.Ethereum.WssURL)
	if err != nil {
		log.Fatal(err)
	}
	contractAddress := common.HexToAddress(args[0])
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		FromBlock: big.NewInt(11043020),
		Topics: [][]common.Hash{
			{
				//common.HexToHash("0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925"),
			},
			{},
			{
				common.HexToHash("0x0000000000000000000000000000000000000000"),
			},
		},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	for _, vLog := range logs {
		evt, err := collector.ParseLog(&vLog, client)
		if err != nil {
			log.Error(err)
			continue
		}
		fmt.Println(evt)
		// collector.LogEvent(&evt)
	}
}
