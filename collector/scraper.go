package collector

import (
	"github.com/iancoleman/strcase"
)

// Protocol is a source token
type Protocol struct {
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	IconURL     string            `json:"icon,omitempty"`
	URL         string            `json:"url,omitempty"`
	Filters     map[string]string `json:"filters,omitempty"`
	Category    string            `json:"category,omitempty"`
}

// ReverseFilters reverse the filters key and value
func (p Protocol) ReverseFilters() map[string]string {
	reversed := make(map[string]string, len(p.Filters))
	for a, n := range p.Filters {
		reversed[strcase.ToCamel(n)] = a
	}
	return reversed
}
