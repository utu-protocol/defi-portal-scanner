package ocean

import (
	"time"
)

type DecentralizedDataObject struct {
	ID         string      `json:"id"`
	NFTAddress string      `json:"nftAddress"`
	Metadata   DDOMetadata `json:"metadata"`
	Datatokens []struct {
		Address   string `json:"address"`
		Name      string `json:"name"`
		Symbol    string `json:"symbol"`
		ServiceID string `json:"serviceId"`
	} `json:"datatokens"`
	Event struct {
		Tx       string `json:"tx"`
		Block    int    `json:"block"`
		From     string `json:"from"`
		Contract string `json:"contract"`
	} `json:"event"`
	Purgatory struct {
		State bool `json:"state"`
	} `json:"purgatory"`
	ChainID int `json:"chainId"`
}

type DDOMetadata struct {
	Author                string    `json:"author"`
	Links                 []string  `json:"links"`
	Tags                  []string  `json:"tags"`
	Categories            []string  `json:"categories"`
	Description           string    `json:"description"`
	Type                  string    `json:"type"`
	Name                  string    `json:"name"`
	DateCreated           time.Time `json:"dateCreated"`
	License               string    `json:"license"`
	AdditionalInformation struct {
		TermsAndConditions bool `json:"termsAndConditions"`
	} `json:"additionalInformation"`
}
