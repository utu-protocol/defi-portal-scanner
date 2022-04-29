package collector

import (
	"github.com/iancoleman/strcase"
	"github.com/utu-crowdsale/defi-portal-scanner/config"
)

// ProtocolsFormat is the format of the protocols.json
type ProtocolsFormat struct {
	OceanPools    []config.Address  `json:"ocean_pools,omitempty"`
	DefiProtocols []Protocol `json:"defi_protocols,omitempty"`
}

// Protocol is a Defi Protocol/Project e.g. Uniswap, Balancer, Sushiswap, whose
// pools (in the Filters field) we listen to for Events, so we can update UTU Trust
// API
type Protocol struct {
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	IconURL     string            `json:"icon,omitempty"`
	URL         string            `json:"url,omitempty"`
	Filters     map[string]string `json:"filters,omitempty"`
	Category    string            `json:"category,omitempty"`
	MainAddress string            `json:"main_address,omitempty"`
}

// ReverseFilters reverse the filters key and value
func (p Protocol) ReverseFilters() map[string]string {
	reversed := make(map[string]string, len(p.Filters))
	for a, n := range p.Filters {
		reversed[strcase.ToCamel(n)] = a
	}
	return reversed
}
