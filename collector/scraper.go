package collector

import (
	"github.com/utu-crowdsale/defi-portal-scanner/utils"
)

// Protocol is a source token
type Protocol struct {
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	IconURL     string            `json:"icon,omitempty"`
	URL         string            `json:"url,omitempty"`
	Filters     map[string]string `json:"filters,omitempty"`
}

// ReverseFilters reverse the filters key and value
func (p Protocol) ReverseFilters() map[string]string {
	reversed := make(map[string]string, len(p.Filters))
	for k, v := range p.Filters {
		reversed[v] = k
	}
	return reversed
}

func getTokens(path string, tokens interface{}) (err error) {

	err = utils.ReadJSON(path, &tokens)
	return

	// url := "https://defimarketcap.io/"
	// c := colly.NewCollector()

	// tokens := []TokenS{}

	// c.OnHTML("tbody tr", func(e *colly.HTMLElement) {

	// 	writer.Write([]string{
	// 		e.ChildText(".cmc-table__column-name"),
	// 		e.ChildText(".cmc-table__cell--sort-by__symbol"),
	// 		e.ChildText(".cmc-table__cell--sort-by__market-cap"),
	// 		e.ChildText(".cmc-table__cell--sort-by__price"),
	// 		e.ChildText(".cmc-table__cell--sort-by__circulating-supply"),
	// 		e.ChildText(".cmc-table__cell--sort-by__volume-24-h"),
	// 		e.ChildText(".cmc-table__cell--sort-by__percent-change-1-h"),
	// 		e.ChildText(".cmc-table__cell--sort-by__percent-change-24-h"),
	// 		e.ChildText(".cmc-table__cell--sort-by__percent-change-7-d"),
	// 	})
	// })

	// c.Visit(url)
}
