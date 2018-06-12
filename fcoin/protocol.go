package fcoin

const (
	OrderTypeLimit = "limit" //限价
)

const (
	OrderSideBuy  = "buy"
	OrderSideSell = "sell"
)

const (
	OrderStateSubmitted       = "submitted"        //已提交
	OrderStatePartialFilled   = "partial_filled"   //部分成交
	OrderStatePartialCanceled = "partial_canceled" //部分成交已撤销
	OrderStateFilled          = "filled"           //完全成交
	OrderStateCanceled        = "canceled"         //已撤销
	OrderStatePendingCancel   = "pending_cancel"   // 撤销已提交
)

type Balance struct {
	Currency  string `json:"currency"`  //
	Available string `json:"available"` //
	Frozen    string `json:"frozen"`    //
	Balance   string `json:"balance"`   //
}

type BalanceReturn struct {
	Status int       `json:"status"` // 请求状态
	Data   []Balance `json:"data"`   // 账户余额

}

type PlaceReturn struct {
	Status int    `json:"status"`
	Data   string `json:"data"`
}

type PlaceRequestParams struct {
	Amount string `json:"amount"` // 限价表示下单数量, 市价买单时表示买多少钱, 市价卖单时表示卖多少币
	Price  string `json:"price"`  // 下单价格, 市价单不传该参数
	Symbol string `json:"symbol"` // 交易对, btcusdt, bccbtc......
	Type   string `json:"type"`   // 订单类型  limit
	Side   string `json:"side"`   //交易方向  buy sell
}

type Order struct {
	ID            string `json:"id"`
	Symbol        string `json:"symbol"`
	Type          string `json:"type"`
	Side          string `json:"side"`
	State         string `json:"state"`
	Amount        string `json:"amount"`
	FilledAmount  string `json:"filled_amount"`
	Price         string `json:"price"`
	ExecutedValue string `json:"executed_value"`
	FillFees      string `json:"fill_fees"`
}
type OrderReturn struct {
	Status int   `json:"status"`
	Data   Order `json:"data"`
}

type OrdersReturn struct {
	Status int     `json:"status"`
	Data   []Order `json:"data"`
}

type MarketDepthData struct {
	Type string    `json:"type"`
	TS   int64     `json:"ts"`
	Seq  string    `json:"seq"`
	Bids []float64 `json:"bids"`
	Asks []float64 `json:"asks"`
}

type MarketDepthReturn struct {
	Status int             `json:"status"`
	Data   MarketDepthData `json:"data"`
}
