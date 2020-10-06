package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// ServicesSchema configure side services
type ServicesSchema struct {
	GlitchtipDsn string `mapstructure:"glitchtip_dsn"`
}

// Schema main configuration for the news room
type Schema struct {
	EthWssURL          string         `mapstructure:"eth_wss_url"`
	EtherscanAPIToken  string         `mapstructure:"etherscan_api_token"`
	TrustEngineURL     string         `mapstructure:"trust_engine_url"`
	DefiSourcesFile    string         `mapstructure:"defi_sources_file"`
	LogOutputFile      string         `mapstructure:"log_output_file"`
	Services           ServicesSchema `mapstructure:"services"`
	RuntimeVersion     string         `mapstructure:"-"`
	RuntimeEnvironment string         `mapstructure:"-"`
	RuntimeName        string         `mapstructure:"-"`
}

// Defaults configure defaults for the configuration
func Defaults() {
	// scheduler defaults
	viper.SetDefault("defi_sources_file", "sources.json")
	viper.SetDefault("log_output_file", "output.json")
}

// Validate a configuration
func Validate(schema *Schema) (err []error) {
	if schema.EthWssURL == "" {
		err = append(err, fmt.Errorf("missing Eth wss URL"))
	}
	if schema.EtherscanAPIToken == "" {
		err = append(err, fmt.Errorf("missing Etherscan API Token"))
	}
	return
}

// Settings general settings
var Settings Schema
