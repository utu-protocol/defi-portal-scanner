package ocean

import (
	"strconv"
)

type Asset struct {
	Pool        *Pool      `json:"pool"`
	Datatoken   *Datatoken `json:"datatoken"`
	PublishedBy string     `json:"published_by"` // this is obtained from pool.controller
	Purgatory   bool       `json:"purgatory"`    // when could this be null? when Aquarius does not have this in the database
	Consumed    uint64     `json:"consumed"`     // Times this data asset was consumed
}

type Pool struct {
	Address          string  `json:"address"`
	Controller       string  `json:"controller"`
	TotalSwapVolume  float64 `json:"total_swap_volume"`
	OceanReserve     float64 `json:"ocean_reserve"`
	DatatokenReserve float64 `json:"datatoken_reserve"`
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

func NewDataToken(address, name, symbol string, orderCount uint64) (dt *Datatoken, err error) {
	return &Datatoken{
		Address:    checksumAddress(address),
		Name:       name,
		Symbol:     symbol,
		OrderCount: orderCount,
	}, nil
}
