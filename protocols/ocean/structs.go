package ocean

import (
	"fmt"
	"strconv"

	"github.com/fatih/structs"
	"github.com/utu-crowdsale/defi-portal-scanner/collector"
)

type Asset struct {
	Pool        *Pool      `json:"pool"`
	Datatoken   *Datatoken `json:"datatoken"`
	PublishedBy string     `json:"published_by"` // this is obtained from pool.controller
	Purgatory   bool       `json:"purgatory"`    // when could this be null? when Aquarius does not have this in the database
	Consumed    uint64     `json:"consumed"`     // Times this data asset was consumed
}

func (a *Asset) Identifier() string {
	return fmt.Sprintf("Asset %s %s", a.Datatoken.Symbol, a.Datatoken.Address)
}

func (a *Asset) toTrustEntity() (te *collector.TrustEntity) {
	te = collector.NewTrustEntity(a.Identifier())

	te.Ids["name"] = a.Datatoken.Name
	te.Ids["symbol"] = a.Datatoken.Symbol
	te.Ids["address_datatoken"] = a.Datatoken.Address
	te.Ids["address_pool"] = a.Pool.Address

	te.Properties = structs.Map(a)
	te.Type = "Asset"

	// These are already represented as other UTU Trust Entity objects, no need
	// to duplicate them as maps here
	delete(te.Properties, "Pool")
	delete(te.Properties, "Datatoken")
	return
}

type Pool struct {
	Address          string  `json:"address"`
	Controller       string  `json:"controller"`
	TotalSwapVolume  float64 `json:"total_swap_volume"`
	OceanReserve     float64 `json:"ocean_reserve"`
	DatatokenReserve float64 `json:"datatoken_reserve"`
}

func (p *Pool) Identifier() string {
	return fmt.Sprintf("Pool %s", p.Address)
}

func (p *Pool) toTrustEntity() (te *collector.TrustEntity) {
	te = collector.NewTrustEntity(fmt.Sprintf("Pool %s", p.Address))
	te.Ids["address"] = p.Address
	te.Properties = structs.Map(p)
	te.Type = "Pool"
	return
}

type DatatokenResponse struct {
	Address    string
	Name       string
	OrderCount string
	Orders     []struct {
		Consumer struct {
			ID string
		}
	}
	Symbol string
}
type PoolGraphQLResponse struct {
	Controller       string
	DatatokenAddress string
	DatatokenReserve string
	ID               string
	OceanReserve     string
	TotalSwapVolume  string
}

func (pgr *PoolGraphQLResponse) toPool() (p *Pool, err error) {
	sv, err := strconv.ParseFloat(pgr.TotalSwapVolume, 64)
	if err != nil {
		return
	}
	or, err := strconv.ParseFloat(pgr.OceanReserve, 64)
	if err != nil {
		return
	}
	dtr, err := strconv.ParseFloat(pgr.DatatokenReserve, 64)
	if err != nil {
		return
	}

	p = &Pool{
		Address:          checksumAddress(pgr.ID),
		Controller:       pgr.Controller,
		TotalSwapVolume:  sv,
		OceanReserve:     or,
		DatatokenReserve: dtr,
	}
	return

}

type Account struct {
	Address     string   `json:"address"`
	AssetsOwned []*Asset `json:"assets_owned"`
	Purgatory   bool     `json:"purgatory"`
}

type Datatoken struct {
	Address    string `json:"address"`     // 0x...
	Name       string `json:"name"`        // Risible Pelican Token
	Symbol     string `json:"symbol"`      // RISPEL-91
	OrderCount uint64 `json:"order_count"` // 1 TokenOrder is one consumption of the asset
}

func (d *Datatoken) Identifier() string {
	return fmt.Sprintf("Datatoken %s (%s)", d.Address, d.Symbol)
}

func (d *Datatoken) toTrustEntity() (te *collector.TrustEntity) {
	te = collector.NewTrustEntity(fmt.Sprintf("Datatoken %s", d.Address))
	te.Ids["address"] = d.Address
	te.Properties = structs.Map(d)
	te.Type = "Datatoken"
	return
}

func NewDataToken(address, name, symbol string, orderCount uint64) (dt *Datatoken) {
	return &Datatoken{
		Address:    checksumAddress(address),
		Name:       name,
		Symbol:     symbol,
		OrderCount: orderCount,
	}
}

func NewDataTokenFromDatatokenResponse(dtr DatatokenResponse) (dt *Datatoken, err error) {
	orderCount, err := strconv.ParseUint(dtr.OrderCount, 10, 64)
	if err != nil {
		return
	}
	return NewDataToken(dtr.Address, dtr.Name, dtr.Symbol, uint64(orderCount)), nil
}
