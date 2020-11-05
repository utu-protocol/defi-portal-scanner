package collector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/utu-crowdsale/defi-portal-scanner/config"
)

// Types definitions
const (
	TypeDefiProtocol = "DeFiProtocol"
	TypeAddress      = "Address"
)

// NewUTUClient create a new utu client
func NewUTUClient(settings config.TrustEngineSchema) *UTUClient {
	return &UTUClient{
		Settings: settings,
		HTTPCli: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// UTUClient trust api client
type UTUClient struct {
	Settings config.TrustEngineSchema
	HTTPCli  *http.Client
}

func (uc UTUClient) postJSON(path string, data interface{}) (err error) {
	if data == nil {
		return
	}
	bin, err := json.Marshal(data)
	if err != nil {
		return
	}
	url := fmt.Sprintf("%s/%s", uc.Settings.URL, path)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bin))
	if err != nil {
		return
	}
	// set the request header
	req.Header.Set("Content-type", "application/json")
	req.Header[uc.Settings.AuthHeder] = []string{uc.Settings.ClientID}
	// execute the request
	rsp, err := uc.HTTPCli.Do(req)
	if err != nil {
		return
	}
	if rsp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(rsp.Body)
		err = fmt.Errorf("server replied with %d: %s\nrequest: %s", rsp.StatusCode, body, bin)
	}
	return
}

// PostEntity post a new entity
func (uc *UTUClient) PostEntity(e *TrustEntity) (err error) {
	err = uc.postJSON("entity", e)
	return
}

// PostRelationship post a new entity
func (uc *UTUClient) PostRelationship(r *TrustRelationship) (err error) {
	err = uc.postJSON("relationship", r)
	return
}

// Apply apply changeset to the utu API
func (uc *UTUClient) Apply(cs *TrustAPIChangeSet) (err error) {
	for _, e := range cs.Entities {
		uc.PostEntity(e)
	}
	for _, r := range cs.Relationship {
		uc.PostRelationship(r)
	}
	return
}

// BuildEntity create an entity that can then be  posted to the trust api
func BuildEntity(name, typ, image string, ids map[string]interface{}, props map[string]string) map[string]interface{} {
	return map[string]interface{}{
		"name":       name,
		"type":       typ,
		"ids":        ids,
		"properties": props,
		"image":      image,
	}
}

// TrustEntity an entity in the trust engine
type TrustEntity struct {
	Type       string                 `json:"type"`
	Name       string                 `json:"name,omitempty"`
	Ids        map[string]string      `json:"ids"`
	Image      string                 `json:"image,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// NewTrustEntity make a new entity
func NewTrustEntity() *TrustEntity {
	return &TrustEntity{
		Ids:        make(map[string]string),
		Properties: make(map[string]interface{}),
	}
}

// TrustRelationship relationship block
type TrustRelationship struct {
	Type           string                 `json:"type"`
	SourceCriteria *TrustEntity           `json:"sourceCriteria"`
	TargetCriteria *TrustEntity           `json:"targetCriteria"`
	Properties     map[string]interface{} `json:"properties"`
}

// NewTrustRelationship make a new relationship
func NewTrustRelationship() *TrustRelationship {
	return &TrustRelationship{
		Properties: make(map[string]interface{}),
	}
}

// TrustAPIChangeSet a changeset to be submitted to the trust api
type TrustAPIChangeSet struct {
	Entities     []*TrustEntity
	Relationship []*TrustRelationship
}

// NewChangeset create a new changeset
func NewChangeset(entitis ...*TrustEntity) *TrustAPIChangeSet {
	return &TrustAPIChangeSet{
		Entities: entitis,
	}
}

// AddRel add a relationship to the changeset
func (cs *TrustAPIChangeSet) AddRel(r *TrustRelationship) {
	cs.Relationship = append(cs.Relationship, r)
}

// AddEntity add a relationship to the changeset
func (cs *TrustAPIChangeSet) AddEntity(e *TrustEntity) {
	cs.Entities = append(cs.Entities, e)
}
