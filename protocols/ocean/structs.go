package ocean

type Asset struct {
	Pool        *Pool      `json:"pool"`
	Datatoken   *Datatoken `json:"datatoken"`
	PublishedBy string     `json:"published_by"` // this is obtained from pool.controller
	Purgatory   bool       `json:"purgatory"`
	Consumed    uint       `json:"consumed"` // Times this data asset was consumed
}

type Pool struct {
	Asset            *Asset `json:"asset"`
	Address          string `json:"address"`
	TotalSwapVolume  uint   `json:"total_swap_volume"`
	OceanReserve     uint   `json:"ocean_reserve"`
	DatatokenReserve uint   `json:"datatoken_reserve"`
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
	OrderCount uint   `json:"order_count"` // 1 TokenOrder is one consumption of the asset
}
