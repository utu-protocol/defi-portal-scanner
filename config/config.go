package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// ServicesSchema configure side services
type ServicesSchema struct {
	GlitchtipDsn string `mapstructure:"glitchtip_dsn"`
}

// ProtocolSchema configuration for protocols
type ProtocolSchema struct {
	Name             string `mapstructure:"name,omitempty"`
	TheGraphEndpoint string `mapstructure:"the_graph_endpoint,omitempty"`
}

// EthereumSchema config for ethereum related resources
type EthereumSchema struct {
	WssURL            string `mapstructure:"node_wss_url"`
	EtherscanAPIToken string `mapstructure:"etherscan_api_token"`
}

// TrustEngineSchema the trust engine client configuration
type TrustEngineSchema struct {
	URL       string `mapstructure:"url"`
	AuthHeder string `mapstructure:"client_id_header"`
	ClientID  string `mapstructure:"client_id"`
	DryRun    bool   `mapstructure:"dry_run"`
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
	Protocols          []ProtocolSchema  `mapstructure:"protocols"`
	RuntimeVersion     string            `mapstructure:"-"`
	RuntimeEnvironment string            `mapstructure:"-"`
	RuntimeName        string            `mapstructure:"-"`
}

// Defaults configure defaults for the configuration
func Defaults() {
	// scheduler defaults
	viper.SetDefault("defi_sources_file", "sources.json")
	viper.SetDefault("log_output_file", "output.json")
	viper.SetDefault("track_topics", []string{"transfer"})
	viper.SetDefault("db_folder", "db")
	// utu api
	viper.SetDefault("utu_trust_api.url", "https://api.ututrust.com")
	viper.SetDefault("utu_trust_api.client_id", "defiPortal")
	viper.SetDefault("utu_trust_api.client_id_header", "UTU-Trust-Api-Client-Id")
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
