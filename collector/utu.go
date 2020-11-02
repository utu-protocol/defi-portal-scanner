package collector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func NewUTUClient(url string) *UTUClient {
	return &UTUClient{
		BaseURL: url,
		HTTPCli: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

type UTUClient struct {
	BaseURL string
	HTTPCli *http.Client
}

func (uc UTUClient) postJson(url string, data interface{}) (err error) {
	bin, err := json.Marshal(data)
	if err != nil {
		return
	}
	rsp, err := uc.HTTPCli.Post(url, "application/json", bytes.NewReader(bin))
	if err != nil {
		return
	}
	if rsp.StatusCode <= http.StatusIMUsed {
		err = fmt.Errorf("server replied with %d", rsp.StatusCode)
	}
	return
}

func (uc *UTUClient) Post(di EthEvent) (err error) {

	entityURL := fmt.Sprintf("%s/postEntity", uc.BaseURL)
	data := map[string]interface{}{
		"type": "client",
		"properties": map[string]interface{}{
			"address": di.ContractAddress,
			"name":    di.ContractName,
			// TODO: also sender and recipient
		},
	}
	err = uc.postJson(entityURL, data)
	return
	//postEntity(type, properties)
	//postRelationship(type, kind, sourceCriteria, targetCriteria, properties, bidirectional)
}
