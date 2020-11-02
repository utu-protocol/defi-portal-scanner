package collector

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/utu-crowdsale/defi-portal-scanner/utils"
)

// Token is a source token
type Token struct {
	Name    string  `json:"name,omitempty"`
	Address string  `json:"address,omitempty"`
	IconURL string  `json:"icon,omitempty"`
	AbiJSON string  `json:"abi,omitempty"`
	Abi     abi.ABI `json:"-"`
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
