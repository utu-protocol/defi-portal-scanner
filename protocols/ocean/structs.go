package ocean

import (
	"fmt"
	"log"
	"strconv"

	"github.com/fatih/structs"
	"github.com/utu-crowdsale/defi-portal-scanner/collector"
)

type Asset struct {
	Name               string     `json:"name"`
	Description        string     `json:"description"`
	DID                string     `json:"did"`
	Datatoken          *Datatoken `json:"datatoken"`
	PublishedBy        string     `json:"published_by"`         // this is obtained from Aquarius DDO
	PublishedByAddress string     `json:"published_by_address"` // this is obtained from pool.controller
	Purgatory          bool       `json:"purgatory"`            // when could this be null? when Aquarius does not have this in the database
	Consumed           uint64     `json:"consumed"`             // Times this data asset was consumed
	Tags               []string   `json:"tags,omitempty"`
	Categories         []string   `json:"categories,omitempty"`
}

func (a *Asset) Identifier() string {
	return fmt.Sprintf("Asset %s %s", a.Datatoken.Symbol, a.Datatoken.Address)
}

func (a *Asset) toTrustEntity() (te *collector.TrustEntity) {
	te = collector.NewTrustEntity(a.Identifier())

	te.Ids["name"] = a.Datatoken.Name
	te.Ids["symbol"] = a.Datatoken.Symbol
	te.Ids["uuid"] = "dt_" + a.Datatoken.Address
	te.Ids["address_datatoken"] = a.Datatoken.Address
	te.Ids["DID"] = a.DID

	te.Properties = structs.Map(a)
	te.Name = a.Name
	te.Type = "Asset"

	// These are already represented as other UTU Trust Entity objects, no need
	// to duplicate them as maps here
	delete(te.Properties, "Datatoken")
	delete(te.Properties, "DID")
	return
}

func (a *Asset) datatokenToTrustRelationship() (tr *collector.TrustRelationship) {
	tr = collector.NewTrustRelationship()
	tr.SourceCriteria = a.toTrustEntity()
	tr.TargetCriteria = a.Datatoken.toTrustEntity()
	tr.Type = "belongsTo"
	return
}

type DatatokenResponse struct {
	Address    string
	Name       string
	OrderCount string
	Orders     []OrderResponse
	NFT        struct {
		Address string
		Creator struct {
			ID string
		}
	}
	Symbol string
}

type OrderResponse struct {
	Tx        string
	Amount    string
	Timestamp uint64
	User      struct {
		ID string
	}
	Price string
}

type Address struct {
	Address               string                  `json:"address"`
	PlaceholderImage      string                  `json:"placeholder_image"`
	Purgatory             bool                    `json:"purgatory"`
	DatatokenInteractions []*DatatokenInteraction `json:"datatoken_interactions"`
}

func (a *Address) toTrustEntity() (te *collector.TrustEntity) {
	te = collector.NewTrustEntity(fmt.Sprintf("Address %s", a.Address))
	te.Ids["address"] = a.Address
	te.Image = a.PlaceholderImage
	te.Properties["purgatory"] = a.Purgatory
	te.Type = "Address"

	return te
}

func (a *Address) datatokenInteractionsToTrustRelationships(datatokensMap map[string]*collector.TrustEntity, log *log.Logger) (tr []*collector.TrustRelationship) {
	for _, dti := range a.DatatokenInteractions {
		t := collector.NewTrustRelationship()
		t.SourceCriteria = a.toTrustEntity()
		x, ok := datatokensMap[dti.AddressDatatoken]
		if !ok {
			log.Printf("%#v mentioned a datatoken %s but I don't know anything about it\n", dti, dti.AddressDatatoken)
			continue
		}
		t.TargetCriteria = x
		t.Type = "interaction"
		t.Properties = structs.Map(dti)
		t.Properties["action"] = "Consumption"
		tr = append(tr, t)
	}
	return tr
}

func (a *Address) Identifier() string {
	return fmt.Sprintf("Address %s", a.Address)
}

func NewAddressFromUserResponse(user string, orders []OrderWrapper, purgatoryMap map[string]string) (a *Address, err error) {
	a = &Address{
		Address:               user,
		Purgatory:             false,
		PlaceholderImage:      fmt.Sprintf("https://via.placeholder.com/150/FFFF00/000000/?text=%s", user),
		DatatokenInteractions: []*DatatokenInteraction{},
	}
	_, ok := purgatoryMap[user]
	if ok {
		a.Purgatory = true
	}

	for _, x := range orders {
		dti := &DatatokenInteraction{
			Address:          user,
			AddressDatatoken: x.Token.ID,
			SymbolDatatoken:  x.Token.Symbol,
			Timestamp:        x.Timestamp,
			TxHash:           x.TxHash,
			Price:            x.Price,
		}
		a.DatatokenInteractions = append(a.DatatokenInteractions, dti)
	}
	return
}

type Datatoken struct {
	Address    string       `json:"address"`     // 0x...
	Name       string       `json:"name"`        // Risible Pelican Token
	Symbol     string       `json:"symbol"`      // RISPEL-91
	OrderCount uint64       `json:"order_count"` // 1 TokenOrder is one consumption of the asset
	NFT        DatatokenNFT `json:"nft"`
	Publisher  string       `json:"string"`
}

type DatatokenNFT struct {
	NFTAddress string `json:"address"`
	Creator    string `json:"creator"`
}

func (d *Datatoken) Identifier() string {
	return fmt.Sprintf("Datatoken %s (%s)", d.Address, d.Symbol)
}

func (d *Datatoken) toTrustEntity() (te *collector.TrustEntity) {
	te = collector.NewTrustEntity(fmt.Sprintf("Datatoken %s", d.Address))
	te.Ids["address"] = d.Address
	te.Properties = structs.Map(d)
	delete(te.Properties, "NFT")
	te.Type = "Datatoken"
	return
}

func NewDataToken(address, name, symbol string, orderCount uint64, publisher string, nftAddress string) (dt *Datatoken) {
	return &Datatoken{
		Address:    address,
		Name:       name,
		Symbol:     symbol,
		OrderCount: orderCount,
		NFT: DatatokenNFT{
			NFTAddress: nftAddress,
			Creator:    publisher,
		},
		Publisher: publisher,
	}
}

func NewDataTokenFromDatatokenResponse(dtr DatatokenResponse) (dt *Datatoken, err error) {
	orderCount, err := strconv.ParseUint(dtr.OrderCount, 10, 64)
	if err != nil {
		return
	}
	return NewDataToken(dtr.Address, dtr.Name, dtr.Symbol, uint64(orderCount), dtr.NFT.Creator.ID, dtr.NFT.Address), nil
}

type DatatokenInteraction struct {
	Address          string `json:"address"`
	AddressDatatoken string `json:"address_datatoken"`
	SymbolDatatoken  string `json:"symbol_datatoken"`
	Timestamp        uint64 `json:"timestamp"`
	TxHash           string `json:"txhash"`
	Price            string `json:"price,omitempty"`
}

type OrderWrapper struct {
	Timestamp uint64
	Amount    string
	TxHash    string
	Token     OrderToken
	Price     string
}

type OrderToken struct {
	ID     string
	Symbol string
	Name   string
}

type OceanConfig struct {
	ChainID      int
	ERC20Address string
	SubgraphURL  string
}

var oceanScanConfigs = []OceanConfig{
	{
		ChainID:      1,
		ERC20Address: "0x967da4048cd07ab37855c090aaf366e4ce1b9f48",
		SubgraphURL:  "https://v4.subgraph.mainnet.oceanprotocol.com/subgraphs/name/oceanprotocol/ocean-subgraph",
	},
	{
		ChainID:      137,
		ERC20Address: "0x282d8efCe846A88B159800bd4130ad77443Fa1A1",
		SubgraphURL:  "https://v4.subgraph.polygon.oceanprotocol.com/subgraphs/name/oceanprotocol/ocean-subgraph",
	},
	{
		ChainID:      56,
		ERC20Address: "0xDCe07662CA8EbC241316a15B611c89711414Dd1a",
		SubgraphURL:  "https://v4.subgraph.bsc.oceanprotocol.com/subgraphs/name/oceanprotocol/ocean-subgraph",
	},
	{
		ChainID:      246,
		ERC20Address: "0x593122AAE80A6Fc3183b2AC0c4ab3336dEbeE528",
		SubgraphURL:  "https://v4.subgraph.energyweb.oceanprotocol.com/subgraphs/name/oceanprotocol/ocean-subgraph",
	},
	{
		ChainID:      1285,
		ERC20Address: "0x99C409E5f62E4bd2AC142f17caFb6810B8F0BAAE",
		SubgraphURL:  "https://v4.subgraph.moonriver.oceanprotocol.com/subgraphs/name/oceanprotocol/ocean-subgraph",
	},
	{
		ChainID:      5,
		ERC20Address: "0xcfdda22c9837ae76e0faa845354f33c62e03653a",
		SubgraphURL:  "https://v4.subgraph.goerli.oceanprotocol.com/subgraphs/name/oceanprotocol/ocean-subgraph",
	},
	{
		ChainID:      80001,
		ERC20Address: "0xd8992Ed72C445c35Cb4A2be468568Ed1079357c8",
		SubgraphURL:  "https://v4.subgraph.mumbai.oceanprotocol.com/subgraphs/name/oceanprotocol/ocean-subgraph",
	},
}
