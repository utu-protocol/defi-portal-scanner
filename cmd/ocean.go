package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/utu-crowdsale/defi-portal-scanner/collector"
	"github.com/utu-crowdsale/defi-portal-scanner/config"
	"github.com/utu-crowdsale/defi-portal-scanner/protocols/ocean"
	"github.com/utu-crowdsale/defi-portal-scanner/utils"
)

var downloadOnly bool

var oceanCmd = &cobra.Command{
	Use:   "ocean",
	Short: "All things OCEAN protocol related",
	Long:  ``,
}

var oceanScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Pull all entities from Ocean Subgraph, Aquarius etc and save to users.json, assets.json",
	Long:  ``,
	RunE:  oceanScan,
}

var oceanPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push to UTU Trust API (expects users.json and assets.json in the current directory)",
	Long:  ``,
	RunE:  oceanPush,
}

var oceanScanPushCmd = &cobra.Command{
	Use:   "scanpush",
	Short: "Combines scan and push... without saving results to intermediate files",
	Long:  ``,
	RunE:  oceanScanPush,
}

func init() {
	oceanCmd.AddCommand(oceanScanCmd)
	oceanCmd.AddCommand(oceanPushCmd)
	oceanCmd.AddCommand(oceanScanPushCmd)
	rootCmd.AddCommand(oceanCmd)
}

func oceanScan(cmd *cobra.Command, args []string) (err error) {
	logger := log.Default()
	assets, users, err := pullDataFromOcean(logger)
	if err != nil {
		return err
	}

	logger.Println("Writing to assets.json")
	if err = utils.WriteJSON("assets.json", true, assets); err != nil {
		return err
	}
	logger.Println("Writing to users.json")
	if err = utils.WriteJSON("users.json", true, users); err != nil {
		return err
	}
	return
}

func oceanPush(cmd *cobra.Command, args []string) (err error) {
	logger := log.Default()
	apiURL, present := os.LookupEnv("APIURL")
	if !present {
		apiURL = "https://stage-api.ututrust.com/core-api"
	}
	apiKey, present := os.LookupEnv("APIKEY")
	if !present {
		return fmt.Errorf("please set the APIKEY environment variable (authorization to the UTU Trust API)")
	}

	var assets []*ocean.Asset
	if err = utils.ReadJSON("assets.json", &assets); err != nil {
		return
	}
	var users []*ocean.Address
	if err = utils.ReadJSON("users.json", &users); err != nil {
		return
	}

	pushToTrustAPI(apiURL, apiKey, assets, users, logger)
	return nil
}

func oceanScanPush(cmd *cobra.Command, args []string) (err error) {
	logger := log.Default()
	apiURL, present := os.LookupEnv("APIURL")
	if !present {
		apiURL = "https://stage-api.ututrust.com/core-api"
	}
	apiKey, present := os.LookupEnv("APIKEY")
	if !present {
		return fmt.Errorf("please set the APIKEY environment variable (authorization to the UTU Trust API)")
	}

	assets, users, err := pullDataFromOcean(logger)
	if err != nil {
		return err
	}

	pushToTrustAPI(apiURL, apiKey, assets, users, logger)
	return nil
}

func pullDataFromOcean(logger *log.Logger) (assets []*ocean.Asset, users []*ocean.Address, err error) {
	logger.Println("Pulling Assets from OCEAN Subgraph")
	assets, err = ocean.PipelineAssets(logger)
	if err != nil {
		return
	}
	logger.Println("Pulling Users from OCEAN Subgraph")
	users, err = ocean.PipelineUsers(logger)
	if err != nil {
		return
	}
	return
}

func pushToTrustAPI(apiURL, apiKey string, assets []*ocean.Asset, users []*ocean.Address, logger *log.Logger) {
	s := &config.TrustEngineSchema{
		URL:           apiURL,
		Authorization: apiKey,
		DryRun:        false,
	}
	utu := collector.NewUTUClient(*s)
	logger.Printf("Posting %d Assets to UTU", len(assets))
	ocean.PostAssetsToUTU(assets, utu, logger)
	logger.Printf("Posting %d Users to UTU", len(users))
	ocean.PostAddressesToUTU(users, assets, utu, logger)
}
