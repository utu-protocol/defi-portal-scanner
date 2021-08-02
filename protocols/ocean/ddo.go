package ocean

import "time"

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
	Service   map[string]interface{} `json:"service"`
	DataToken string                 `json:"dataToken"`
	Created   time.Time              `json:"created"`
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
