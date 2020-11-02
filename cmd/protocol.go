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
	"sync"

	"github.com/spf13/cobra"
	"github.com/utu-crowdsale/defi-portal-scanner/protocols/uniswap"
)

var (
	pwg       sync.WaitGroup
	uniswapOn bool
)

// protocolCmd represents the protocol command
var protocolCmd = &cobra.Command{
	Use:   "protocol",
	Short: "Enable collecting data from a protocol",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		if uniswapOn {
			theGraphEndpoint := "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v2"
			go uniswap.Start(theGraphEndpoint)
			pwg.Add(1)
		}
		pwg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(protocolCmd)
	//
	protocolCmd.Flags().BoolVar(&uniswapOn, "uniswap", false, "enable uniswap protocol")
}
