package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// ServicesSchema configure side services
type ServicesSchema struct {
	GlitchtipDsn string `mapstructure:"glitchtip_dsn"`
}

// EthereumSchema config for ethereum related resources
type EthereumSchema struct {
	WssURL            string `mapstructure:"node_wss_url"`
	EtherscanAPIToken string `mapstructure:"etherscan_api_token"`
}

// TrustEngineSchema the trust engine client configuration
type TrustEngineSchema struct {
	URL           string `mapstructure:"url"`
	Authorization string `mapstructure:"authorization"`
	DryRun        bool   `mapstructure:"dry_run"`
}

// ServerSchema the schema for server
type ServerSchema struct {
	ListenAddress string `mapstructure:"listen_address"`
}

// Schema main configuration for the news room
type Schema struct {
	Ethereum           EthereumSchema    `mapstructure:"eth"`
	UTUTrustAPI        TrustEngineSchema `mapstructure:"utu_trust_api"`
	DefiSourcesFile    string            `mapstructure:"defi_sources_file"`
	TrackTopics        []string          `mapstructure:"track_topics"`
	DbFolder           string            `mapstructure:"db_folder"`
	LogOutputFile      string            `mapstructure:"log_output_file"`
	Services           ServicesSchema    `mapstructure:"services"`
	Server             ServerSchema      `mapstructure:"server"`
	RuntimeVersion     string            `mapstructure:"-"`
	RuntimeEnvironment string            `mapstructure:"-"`
	RuntimeName        string            `mapstructure:"-"`
}

// Defaults configure defaults for the configuration
func Defaults() {
	// scheduler defaults
	viper.SetDefault("defi_sources_file", "protocols.json")
	viper.SetDefault("log_output_file", "output.json")
	viper.SetDefault("track_topics", []string{"transfer"})
	viper.SetDefault("db_folder", "db")
	// utu api
	viper.SetDefault("utu_trust_api.url", "https://api.ututrust.com")
	viper.SetDefault("utu_trust_api.client_id", "defiPortal")
	viper.SetDefault("utu_trust_api.client_id_header", "UTU-Trust-Api-Client-Id")
	// server
	viper.SetDefault("server.listen_address", ":2011")
}

// Validate a configuration
func Validate(schema *Schema) (err []error) {
	if schema.Ethereum.WssURL == "" {
		err = append(err, fmt.Errorf("missing Eth wss URL"))
	}
	if schema.Ethereum.EtherscanAPIToken == "" {
		err = append(err, fmt.Errorf("missing Etherscan API Token"))
	}
	return
}

// Settings general settings
var Settings Schema
