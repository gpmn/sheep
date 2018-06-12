package proto

const (
	AccountBalanceTypeTrade  = "trade"
	AccountBalanceTypeFrozen = "frozen"
)

type AccountBalance struct {
	Currency string `json:"currency"` // 币种 btc eth etc
	Balance  string `json:"balance"`  // 结余
	Type     string `json:"type"`     // 类型, trade: 交易余额, frozen: 冻结余额
}

const (
	OrderPlaceTypeBuyMarket  = "buy-market"  //市价买
	OrderPlaceTypeSellMarket = "sell-market" //市价卖
	OrderPlaceTypeBuyLimit   = "buy-limit"   //限价买
	OrderPlaceTypeSellLimit  = "sell-limit"  //限价卖
)

type OrderPlaceParams struct {
	Price           float64 `json:"price"`
	Amount          float64 `json:"amount"`
	BaseCurrencyID  string  `json:"base_currency_id"`
	QuoteCurrencyID string  `json:"quote_currency_id"`
	Type            string  `json:"type"`
}

type OrderPlaceReturn struct {
	OrderID string `json:"order_id"`
}

type OrderCancelParams struct {
	OrderID         string `json:"order_id"`
	BaseCurrencyID  string `json:"base_currency_id"`  //OKEX 必填
	QuoteCurrencyID string `json:"quote_currency_id"` //OKEX 必填
}

type OrderInfoParams struct {
	OrderID         string `json:"order_id"`
	BaseCurrencyID  string `json:"base_currency_id"`  //OKEX 必填
	QuoteCurrencyID string `json:"quote_currency_id"` //OKEX 必填
}

const (
	OrderStateFilled        = "filled"         //完全成交
	OrderStateSubmitted     = "submitted"      //已提交
	OrderStateCanceled      = "canceled"       //已撤销
	OrderStatePartialFilled = "partial-filled" //部分成交
)

type Order struct {
	ID          string  `json:"id"`
	Symbol      string  `json:"symbol"`
	State       string  `json:"state"`
	Amount      float64 `json:"amount"`
	FieldAmount float64 `json:"field-amount"`
	Price       float64 `json:"price"`
	Type        string  `json:"type"`
}

type OrdersParams struct {
	Symbol          string `json:"symbol"`
	States          string `json:"states"`
	BaseCurrencyID  string `json:"base_currency_id"`  //OKEX 必填
	QuoteCurrencyID string `json:"quote_currency_id"` //OKEX 必填
	Status          string `json:"status"`            //OKEX 查询状态 0：未完成的订单 1：已经完成的订单 （最近两天的数据）
	CurrentPage     string `json:"current_page"`      //OKEX 当前页数
	PageLength      string `json:"page_length"`       //OKEX 每页数据条数，最多不超过200
}

type MarketDepthParams struct {
	Symbol string `json:"symbol"`
	Level  string `json:"level"`
}
