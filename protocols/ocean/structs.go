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

type User struct {
	Address               string                  `json:"address"`
	Purgatory             bool                    `json:"purgatory"`
	DatatokenInteractions []*DatatokenInteraction `json:"datatoken_interactions"`
	PoolInteractions      []*PoolInteraction      `json:"pool_interactions"`
}

func NewUserFromUserResponse(ur UserResponse, purgatoryMap map[string]string) (u *User, err error) {
	/*
		{
			"id": "0x006d0f31a00e1f9c017ab039e9d0ba699433a28c",
			"orders": [
				{
				"amount": "1",
				"datatokenId": {
					"id": "0xfcb47f5781f14ed7e032bd395113b84c897aa23f",
					"name": "Trenchant Pelican Token",
					"symbol": "TREPEL-36"
				},
				"timestamp": 1629082751
				}
			],
			"poolTransactions": [
				{
				"event": "swap",
				"poolAddressStr": "0xa94a4ed3b3414bb2468e5c200d68e56d4ce180f9",
				"sharesTransferAmount": "0",
				"timestamp": 1605717888
				},
			]
		}
	*/

	u = &User{
		Address:               ur.ID,
		Purgatory:             false,
		DatatokenInteractions: []*DatatokenInteraction{},
		PoolInteractions:      []*PoolInteraction{},
	}
	_, ok := purgatoryMap[ur.ID]
	if ok {
		u.Purgatory = true
	}

	for _, x := range ur.Orders {
		dti := &DatatokenInteraction{
			AddressUser:      ur.ID,
			AddressDatatoken: x.DatatokenID.ID,
			SymbolDatatoken:  x.DatatokenID.Symbol,
			Timestamp:        x.Timestamp,
		}
		u.DatatokenInteractions = append(u.DatatokenInteractions, dti)
	}
	for _, x := range ur.PoolTransactions {
		p := &PoolInteraction{
			AddressUser: ur.ID,
			AddressPool: x.PoolAddress,
			Event:       x.Event,
			Timestamp:   x.Timestamp,
		}
		u.PoolInteractions = append(u.PoolInteractions, p)
	}
	return
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

type DatatokenInteraction struct {
	AddressUser      string `json:"address_user"`
	AddressDatatoken string `json:"address_datatoken"`
	SymbolDatatoken  string `json:"symbol_datatoken"`
	Timestamp        uint64 `json:"timestamp"`
}

type PoolInteraction struct {
	AddressUser string `json:"address_user"`
	AddressPool string `json:"address_pool"`
	Event       string `json:"event"`
	Timestamp   uint64 `json:"timestamp"`
}

type UserResponse struct {
	ID     string `json:"id"`
	Orders []struct {
		Timestamp   uint64 `json:"timestamp"`
		Amount      string `json:"amount"`
		DatatokenID struct {
			ID     string `json:"id"`
			Name   string `json:"name"`
			Symbol string `json:"symbol"`
		} `json:"datatokenId"`
	} `json:"orders"`
	PoolTransactions []struct {
		Event                string `json:"event"`
		PoolAddress          string `json:"poolAddressStr"`
		SharesTransferAmount string `json:"sharesTransferAmount"`
		Timestamp            uint64 `json:"timestamp"`
	} `json:"poolTransactions"`
}
