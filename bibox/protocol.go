package bibox

type Ask struct {
	Price  string `json:"price"`
	Volume string `json:"volume"`
}

type Bid struct {
	Price  string `json:"price"`
	Volume string `json:"volume"`
}

type GetMarketDepthRspResult struct {
	Pair       string `json:"pair"`
	UpdateTime int64  `json:"update_time"`
	Asks       []Ask  `json:"asks"`
	Bids       []Bid  `json:"bids"`
}

type GetMarketDepthRsp struct {
	Cmd    string                  `json:"cmd"`
	Result GetMarketDepthRspResult `json:"result"`
}

type Cmd struct {
	Cmd  string            `json:"cmd"`
	Body map[string]string `json:"body"`
}

type ApiKeyReq struct {
	Cmds   string `json:"cmds"`
	Apikey string `json:"apikey"`
	Sign   string `json:"sign"`
}

type GetAccountBalanceRspAsset struct {
	CoinSymbol string `json:"coin_symbol"`
	Balance    string `json:"balance"`
	Freeze     string `json:"freeze"`
	USDValue   string `json:"USDValue"`
}

type GetAccountBalanceRspResult struct {
	TotalBTC   string                      `json:"total_btc"`
	TotalCNY   string                      `json:"total_cny"`
	TotalUSD   string                      `json:"total_usd"`
	AssetsList []GetAccountBalanceRspAsset `json:"assets_list"`
}

type GetAccountBalanceRsp struct {
	Result []struct {
		Result GetAccountBalanceRspResult `json:"result"`
	} `json:"result"`
	Cmd string `json:"cmd"`
}

type OrderPlaceRspResult struct {
	Result int    `json:"result"`
	Cmd    string `json:"cmd"`
}

type OrderPlaceRsp struct {
	Result []OrderPlaceRspResult `json:"result"`
}

type OrderCancelRsp struct {
	Result string `json:"result"`
	Cmd    string `json:"cmd"`
}

type OrderPendingListRspResultItem struct {
	ID             int64  `json:"id"`
	CreatedAt      int64  `json:"createdAt"`
	AccountType    int    `json:"account_type"`
	CoinSymbol     string `json:"coin_symbol"`
	CurrencySymbol string `json:"currency_symbol"`
	OrderSide      int    `json:"order_side"`
	OrderType      int    `json:"order_type"`
	Price          string `json:"price"`
	Amount         string `json:"amount"`
	Money          string `json:"money"`
	DealPrice      string `json:"deal_price"`
	DealAmount     string `json:"deal_amount"`
	Unexecuted     string `json:"unexecuted"`
	DealPercent    string `json:"deal_percent"`
	Status         int    `json:"status"`
}

type OrderPendingListRspResult struct {
	Result struct {
		Count int                             `json:"count"`
		Page  int                             `json:"page"`
		Items []OrderPendingListRspResultItem `json:"items"`
	} `json:"result"`
	Cmd string `json:"cmd"`
}

type OrderPendingListRsp struct {
	Result []OrderPendingListRspResult `json:"result"`
}

type OrderInfoRspResult struct {
	ID             int64  `json:"id"`
	CreatedAt      int64  `json:"createdAt"`
	AccountType    int    `json:"account_type"`
	Pair           string `json:"pair"`
	CoinSymbol     string `json:"coin_symbol"`
	CurrencySymbol string `json:"currency_symbol"`
	OrderSide      int    `json:"order_side"`
	OrderType      int    `json:"order_type"`
	Price          string `json:"price"`
	Amount         string `json:"amount"`
	Money          string `json:"money"`
	DealPrice      string `json:"deal_price"`
	DealAmount     string `json:"deal_amount"`
	DealMoney      string `json:"deal_money"`
	DealPercent    string `json:"deal_percent"`
	Status         int    `json:"status"`
	Cmd            string `json:"cmd"`
}

type OrderInfoRsp struct {
	Result []struct {
		Result OrderInfoRspResult `json:"result"`
		Cmd    string             `json:"cmd"`
	} `json:"result"`
}

type GetOrderHistoryListRspResult struct {
	ID             uint   `json:"id"`
	CreatedAt      int64  `json:"createdAt"`
	AccountType    int    `json:"account_type"`
	CoinSymbol     string `json:"coin_symbol"`
	CurrencySymbol string `json:"currency_symbol"`
	OrderSide      int    `json:"order_side"`
	OrderType      int    `json:"order_type"`
	Price          string `json:"price"`
	Amount         string `json:"amount"`
	Money          string `json:"money"`
	Fee            string `json:"fee"`
}

type GetOrderHistoryListRsp struct {
	Result []struct {
		Result struct {
			Count int                            `json:"count"`
			Page  int                            `json:"page"`
			Items []GetOrderHistoryListRspResult `json:"items"`
		} `json:"result"`
		Cmd string `json:"cmd"`
	} `json:"result"`
}
