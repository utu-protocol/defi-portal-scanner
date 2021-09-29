package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/utu-crowdsale/defi-portal-scanner/collector"
	"github.com/utu-crowdsale/defi-portal-scanner/config"
	"github.com/utu-crowdsale/defi-portal-scanner/protocols/ocean"
)

var oceanCmd = &cobra.Command{
	Use:   "ocean",
	Short: "All things OCEAN protocol related",
	Long:  ``,
}

var scanpushCmd = &cobra.Command{
	Use:   "scanpush",
	Short: "Pull all entities from Ocean Subgraph, Aquarius etc and push to UTU Trust API",
	Long:  ``,
	RunE:  scanPush,
}

func init() {
	oceanCmd.AddCommand(scanpushCmd)
	rootCmd.AddCommand(oceanCmd)
}

func scanPush(cmd *cobra.Command, args []string) (err error) {
	logger := log.Default()
	apiURL, present := os.LookupEnv("APIURL")
	if !present {
		apiURL = "https://stage-api.ututrust.com/core-api"
	}
	apiKey, present := os.LookupEnv("APIKEY")
	if !present {
		return fmt.Errorf("please set the APIKEY environment variable (authorization to the UTU Trust API)")
	}

	logger.Println("Pulling Assets from OCEAN Subgraph")
	assets, err := ocean.PipelineAssets(logger)
	if err != nil {
		return err
	}
	logger.Println("Pulling Users from OCEAN Subgraph")
	users, err := ocean.PipelineUsers(logger)
	if err != nil {
		return err
	}

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

	return nil
}
