package ocean

import (
	"encoding/json"
	"fmt"
	"time"
)

type DecentralizedDataObject struct {
	Context   string `json:"@context"`
	ID        string `json:"id"`
	PublicKey []struct {
		ID    string `json:"id"`
		Type  string `json:"type"`
		Owner string `json:"owner"`
	} `json:"publicKey"`
	Authentication []struct {
		Type      string `json:"type"`
		PublicKey string `json:"publicKey"`
	} `json:"authentication"`
	Service   []interface{} `json:"service"`
	DataToken string        `json:"dataToken"`
	Created   time.Time     `json:"created"`
	Proof     struct {
		Created        time.Time `json:"created"`
		Creator        string    `json:"creator"`
		Type           string    `json:"type"`
		SignatureValue string    `json:"signatureValue"`
	} `json:"proof"`
	DataTokenInfo struct {
		Address  string  `json:"address"`
		Name     string  `json:"name"`
		Symbol   string  `json:"symbol"`
		Decimals int     `json:"decimals"`
		Cap      float64 `json:"cap"`
	} `json:"dataTokenInfo"`
	Updated         time.Time     `json:"updated"`
	AccessWhiteList []interface{} `json:"accessWhiteList"`
	Event           struct {
		Txid     string `json:"txid"`
		BlockNo  int    `json:"blockNo"`
		From     string `json:"from"`
		Contract string `json:"contract"`
		Update   bool   `json:"update"`
	} `json:"event"`
	IsInPurgatory string `json:"isInPurgatory"`
	ChainID       int    `json:"chainId"`
}

func (ddo *DecentralizedDataObject) GetNameAuthorMetadata() (name, author, description string, tags, categories []string, err error) {
	for _, obj := range ddo.Service {
		if obj.(map[string]interface{})["type"] == "metadata" {
			var ddomJson []byte
			ddomJson, err = json.Marshal(obj)
			if err != nil {
				return
			}
			ddom := new(DDOMetadata)
			err = json.Unmarshal(ddomJson, ddom)
			if err != nil {
				return
			}

			name = ddom.Attributes.Main.Name
			author = ddom.Attributes.Main.Author
			description = ddom.Attributes.AdditionalInformation.Description
			tags = ddom.Attributes.AdditionalInformation.Tags
			categories = ddom.Attributes.AdditionalInformation.Categories

			return name, author, description, tags, categories, nil
		}
	}
	return "", "", "", nil, nil, fmt.Errorf("DecentralizedDataObject.GetNameDescription() couldn't find a name and description")
}

type DDOMetadata struct {
	Type       string `json:"type"`
	Attributes struct {
		Curation struct {
			Rating   int  `json:"rating"`
			NumVotes int  `json:"numVotes"`
			IsListed bool `json:"isListed"`
		} `json:"curation"`
		Main struct {
			Type        string    `json:"type"`
			Name        string    `json:"name"`
			DateCreated time.Time `json:"dateCreated"`
			Author      string    `json:"author"`
			License     string    `json:"license"`
			Files       []struct {
				ContentLength string `json:"contentLength"`
				ContentType   string `json:"contentType"`
				Index         int    `json:"index"`
			} `json:"files"`
			DatePublished time.Time `json:"datePublished"`
		} `json:"main"`
		AdditionalInformation struct {
			Description string   `json:"description"`
			Tags        []string `json:"tags"`
			Categories  []string `json:"categories"`
			Links       []struct {
				ContentLength string `json:"contentLength"`
				ContentType   string `json:"contentType"`
				URL           string `json:"url"`
			} `json:"links"`
			TermsAndConditions bool `json:"termsAndConditions"`
		} `json:"additionalInformation"`
		EncryptedFiles string `json:"encryptedFiles"`
	} `json:"attributes"`
	Index int `json:"index"`
}
