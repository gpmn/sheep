package huobi

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"

	"fmt"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/gpmn/sheep/consts"
	"github.com/gpmn/sheep/proto"
)

// LoanOrder :
type LoanOrder struct {
	AccountID       int     `json:"account-id"`
	AccruedAt       int     `json:"accrued-at"`
	CreatedAt       int     `json:"created-at"`
	Currency        string  `json:"currency"`
	ID              int     `json:"id"`
	InterestAmount  string  `json:"interest-amount"`
	InterestBalance float64 `json:"interest-balance,string"`
	InterestRate    string  `json:"interest-rate"`
	LoanAmount      string  `json:"loan-amount"`
	LoanBalance     float64 `json:"loan-balance,string"`
	State           string  `json:"state"`
	Symbol          string  `json:"symbol"`
	UserID          int     `json:"user-id"`
}

// LoanOrderResp :
type LoanOrderResp struct {
	Orders []LoanOrder `json:"data"`
	Status string      `json:"status"`
}

// MarginBalanceItem :
type MarginBalanceItem struct {
	Balance  float64 `json:"balance,string"`
	Currency string  `json:"currency"`
	Type     string  `json:"type"`
}

// FindMBItem :
func (mb *MarginBalance) FindMBItem(currency, tp string) *MarginBalanceItem {
	if mb.State != "working" {
		log.Printf("MarginBalance.FindMBItem - mb.State %s invalid", mb.State)
		return nil
	}

	for idx := range mb.List {
		item := &mb.List[idx]
		if item.Currency == currency && item.Type == tp {
			return item
		}
	}
	return nil
}

// MarginBalance :
type MarginBalance struct {
	FlPrice  string              `json:"fl-price"`
	FlType   string              `json:"fl-type"`
	ID       int                 `json:"id"`
	List     []MarginBalanceItem `json:"list"`
	RiskRate string              `json:"risk-rate"`
	State    string              `json:"state"`
	Symbol   string              `json:"symbol"`
	Type     string              `json:"type"`
}

// MarginBalanceResp :
type MarginBalanceResp struct {
	Balances []MarginBalance `json:"data"`
	Status   string          `json:"status"`
}

// TickData :
type TickData struct {
	Amount    float64 `json:"amount"`
	Direction string  `json:"direction"`
	Price     float64 `json:"price"`
	TS        int64   `json:"ts"`
}

// MarketTradeDetail :
type MarketTradeDetail struct {
	Ch   string `json:"ch"`
	Tick struct {
		Data []TickData `json:"data"`
	} `json:"tick"`
}

func (m *MarketTradeDetail) String() string {
	return fmt.Sprintln(m.Ch, "实时价格推送  价格:", m.Tick.Data[0].Price, " 数量:", m.Tick.Data[0].Amount, " 买卖：", m.Tick.Data[0].Direction)
}

type MarketDepth struct {
	Ch   string `json:"ch"`
	Tick struct {
		Asks [][]float64 `json:"asks"`
		Bids [][]float64 `json:"bids"`
		TS   int64       `json:"ts"`
	} `json:"tick"`
}

type Account struct {
	ID     int64
	Type   string
	State  string
	UserID int64
}

type Huobi struct {
	accessKey       string
	secretKey       string
	tradeAccount    Account
	market          *Market
	depthListener   DepthlListener
	detailListener  DetailListener
	klineUpListener KLineUpListener
	orderListener   OrderListener
}

func (h *Huobi) OpenWebsocket() error {
	var err error
	h.market, err = NewMarket()
	if err != nil {
		return err
	}

	go h.market.Loop()
	return nil
}

func (h *Huobi) CloseWebsocket() error {
	return h.market.Close()
}

func (h *Huobi) GetExchangeType() string {
	return consts.ExchangeTypeHuobi
}

type KLine struct {
	ID     int     `json:"id"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	Low    float64 `json:"low"`
	High   float64 `json:"high"`
	Amount float64 `json:"amount"`
	Vol    float64 `json:"vol"`
	Count  int     `json:"count"`
}

type KLineUpdate struct {
	Ch    string `json:"ch"`
	Ts    int64  `json:"ts"`
	Kline KLine  `json:"tick"`
}

type RespGetKLines struct {
	Status string  `json:"status"`
	Ch     string  `json:"ch"`
	Ts     int64   `json:"ts"`
	KLines []KLine `json:"data"`
}

// 查询当前用户的K线数据
func (h *Huobi) GetKLines(symbol, period string, size int) (kl RespGetKLines, err error) {
	strRequest := "/market/history/kline"
	param := make(map[string]string)
	param["symbol"] = symbol
	param["period"] = period
	param["size"] = strconv.Itoa(size)
	jsonReturn, err := apiKeyGet(param, strRequest, h.accessKey, h.secretKey)
	if nil != err {
		log.Printf("Huobi.GetKLines - apiKeyGet failed : %v, content : %s", err, jsonReturn)
		return kl, err
	}
	err = json.Unmarshal([]byte(jsonReturn), &kl)
	if nil != err {
		log.Printf("Huobi.GetKLines - failed : %v, content : %s", err, jsonReturn)
	}
	return kl, err
}

// GetAccounts : 查询当前用户的所有账户, 根据包含的私钥查询
// return: AccountsReturn对象
func (h *Huobi) GetAccounts() (AccountsReturn, error) {
	accountsReturn := AccountsReturn{}

	strRequest := "/v1/account/accounts"
	buf, err := apiKeyGet(make(map[string]string), strRequest, h.accessKey, h.secretKey)
	if nil != err {
		log.Printf("Huobi.GetAccounts - apiKeyGet failed : %v, content : %s", err, buf)
		return accountsReturn, err
	}
	json.Unmarshal([]byte(buf), &accountsReturn)
	return accountsReturn, nil
}

// GetAccountBalance : 根据账户ID查询账户余额
// return: BalanceReturn对象
func (h *Huobi) GetAccountBalance() ([]proto.AccountBalance, error) {
	balanceReturn := BalanceReturn{}
	strRequest := fmt.Sprintf("/v1/account/accounts/%d/balance", h.tradeAccount.ID)
	jsonBanlanceReturn, err := apiKeyGet(make(map[string]string), strRequest, h.accessKey, h.secretKey)
	if nil != err {
		log.Printf("Huobi.GetAccountBalance - apiKeyGet failed : %v", err)
		return nil, err
	}
	if err = json.Unmarshal([]byte(jsonBanlanceReturn), &balanceReturn); err != nil {
		log.Printf("Huobi.GetAccountBalance - json.Unmarshal failed : %v", err)
		return nil, err
	}
	if balanceReturn.Status != "ok" {
		log.Printf("Huobi.GetAccountBalance - json.Unmarshal errmsg : %s", balanceReturn.ErrMsg)
		return nil, fmt.Errorf(balanceReturn.ErrMsg)
	}

	var res []proto.AccountBalance
	for _, blance := range balanceReturn.Data.List {
		var item proto.AccountBalance
		item.Currency = blance.Currency
		item.Balance = blance.Balance
		item.Type = blance.Type

		res = append(res, item)
	}

	return res, nil
}

// GetMarginBalances : 根据账户ID查询账户余额
func (h *Huobi) GetMarginBalances(baseSym string) ([]MarginBalance, error) {
	//balanceReturn := BalanceReturn{}
	strRequest := fmt.Sprintf("/v1/margin/accounts/balance")
	args := make(map[string]string)
	if baseSym != "" {
		args["symbol"] = baseSym
	}
	buf, err := apiKeyGet(args, strRequest, h.accessKey, h.secretKey)
	if nil != err {
		log.Printf("Huobi.GetMarginBalances - apiKeyGet failed : %v", err)
		return nil, err
	}
	var mbr MarginBalanceResp
	if err = json.Unmarshal([]byte(buf), &mbr); nil != err {
		log.Printf("Huobi.GetMarginBalances - json.Unmarshal failed %v", err)
		return nil, err
	}
	if mbr.Status != "ok" {
		log.Printf("Huobi.GetMarginBalances - response status %s, invalid", mbr.Status)
		return nil, fmt.Errorf("response status wrong : " + mbr.Status)
	}

	return mbr.Balances, nil
}

// OrderPlace :下单
// placeRequestParams: 下单信息
// return: OrderID
func (h *Huobi) OrderPlace(params *proto.OrderPlaceParams) (*proto.OrderPlaceReturn, error) {
	placeReturn := PlaceReturn{}
	var placeRequestParams PlaceRequestParams
	placeRequestParams.AccountID = strconv.FormatInt(h.tradeAccount.ID, 10)
	fmtStr := fmt.Sprintf("%%.%df", params.AmountPrecision)
	placeRequestParams.Amount = fmt.Sprintf(fmtStr, params.Amount) //strconv.FormatFloat(params.Amount, 'f', -1, 64)
	fmtStr = fmt.Sprintf("%%.%df", params.PricePrecision)
	placeRequestParams.Price = fmt.Sprintf(fmtStr, params.Price)
	placeRequestParams.Source = "api"
	placeRequestParams.Symbol = strings.ToLower(params.BaseCurrencyID) + strings.ToLower(params.QuoteCurrencyID)
	placeRequestParams.Type = params.Type

	mapParams := make(map[string]string)
	mapParams["account-id"] = placeRequestParams.AccountID
	mapParams["amount"] = placeRequestParams.Amount
	if 0 < len(placeRequestParams.Price) {
		mapParams["price"] = placeRequestParams.Price
	}
	if 0 < len(placeRequestParams.Source) {
		mapParams["source"] = placeRequestParams.Source
	}
	mapParams["symbol"] = placeRequestParams.Symbol
	mapParams["type"] = placeRequestParams.Type

	strRequest := "/v1/order/orders/place"
	buf, err := apiKeyPost(mapParams, strRequest, h.accessKey, h.secretKey)
	if nil != err {
		log.Printf("Huobi.OrderPlace - apiKeyPost failed : %v", err)
		return nil, err
	}
	json.Unmarshal([]byte(buf), &placeReturn)

	if placeReturn.Status != "ok" {
		log.Printf("Huobi.OrderPlace - status wrong : %v", placeReturn)
		return nil, errors.New(placeReturn.ErrMsg)
	}

	var ret proto.OrderPlaceReturn
	ret.OrderID = placeReturn.Data

	return &ret, nil

}

// OrderCancel : 申请撤销一个订单请求
// strOrderID: 订单ID
// return: PlaceReturn对象
func (h *Huobi) OrderCancel(params *proto.OrderCancelParams) error {
	placeReturn := PlaceReturn{}

	strRequest := fmt.Sprintf("/v1/order/orders/%s/submitcancel", params.OrderID)
	buf, err := apiKeyPost(make(map[string]string), strRequest, h.accessKey, h.secretKey)
	if nil != err {
		log.Printf("Huobi.OrderCancel - apiKeyPost failed : %v", err)
		return err
	}
	json.Unmarshal([]byte(buf), &placeReturn)

	if placeReturn.Status != "ok" {
		return errors.New(placeReturn.ErrMsg)
	}

	return nil
}

// 查询订单详情
// strOrderID: 订单ID
// return: OrderReturn对象
func (h *Huobi) GetOrderInfo(params *proto.OrderInfoParams) (*proto.Order, error) {
	orderReturn := OrderReturn{}

	strRequest := fmt.Sprintf("/v1/order/orders/%s", params.OrderID)
	jsonPlaceReturn, err := apiKeyGet(make(map[string]string), strRequest, h.accessKey, h.secretKey)
	if nil != err {
		log.Printf("Huobi.GetOrderInfo - apiKeyGet failed : %v", err)
		return nil, err
	}
	json.Unmarshal([]byte(jsonPlaceReturn), &orderReturn)

	if orderReturn.Status != "ok" {
		return nil, errors.New(orderReturn.ErrMsg)
	}

	var ret proto.Order
	ret.Price, _ = strconv.ParseFloat(orderReturn.Data.Price, 64)
	ret.ID = strconv.FormatInt(orderReturn.Data.ID, 10)
	ret.Symbol = orderReturn.Data.Symbol
	ret.State = orderReturn.Data.State
	ret.FieldAmount, _ = strconv.ParseFloat(orderReturn.Data.FieldAmount, 64)
	ret.Type = orderReturn.Data.Type
	ret.Amount, _ = strconv.ParseFloat(orderReturn.Data.Amount, 64)

	return &ret, nil

}

// GetOrders :
func (h *Huobi) GetOrders(params *proto.OrdersParams) ([]proto.Order, error) {
	ordersReturn := OrdersReturn{}

	jsonP, _ := json.Marshal(params)

	var paramMap = make(map[string]string)
	json.Unmarshal(jsonP, &paramMap)

	strRequest := "/v1/order/orders"
	jsonRet, err := apiKeyGet(paramMap, strRequest, h.accessKey, h.secretKey)
	if nil != err {
		log.Printf("Huobi.GetOrders - apiKeyGet failed : %v", err)
		return nil, err
	}

	json.Unmarshal([]byte(jsonRet), &ordersReturn)
	if ordersReturn.Status != "ok" {
		return nil, errors.New(ordersReturn.ErrMsg)
	}

	var ret []proto.Order
	for _, cell := range ordersReturn.Data {
		var item proto.Order
		item.Price, _ = strconv.ParseFloat(cell.Price, 64)
		item.ID = strconv.FormatInt(cell.ID, 10)
		item.Symbol = cell.Symbol
		item.State = cell.State
		item.FieldAmount, _ = strconv.ParseFloat(cell.FieldAmount, 64)
		item.Type = cell.Type
		item.Amount, _ = strconv.ParseFloat(cell.Amount, 64)
		item.CreatedSec = cell.CreatedSec

		ret = append(ret, item)
	}

	return ret, nil
}

// GetOpenOrders :
func (h *Huobi) GetOpenOrders(params *proto.OrdersParams) ([]proto.Order, error) {
	oor := OpenOrdersReturn{}

	jsonP, _ := json.Marshal(params)

	var paramMap = make(map[string]string)
	json.Unmarshal(jsonP, &paramMap)

	strRequest := "/v1/order/openOrders"
	jsonRet, err := apiKeyGet(paramMap, strRequest, h.accessKey, h.secretKey)
	if nil != err {
		log.Printf("Huobi.GetOrders - apiKeyGet failed : %v", err)
		return nil, err
	}

	json.Unmarshal([]byte(jsonRet), &oor)
	if oor.Status != "ok" {
		return nil, errors.New(oor.Status)
	}

	var ret []proto.Order
	for idx := range oor.Data {
		cell := &oor.Data[idx]
		var item proto.Order
		item.Price = cell.Price
		item.ID = fmt.Sprintf("%d", cell.ID)
		item.Symbol = cell.Symbol
		item.State = cell.State
		item.FieldAmount = cell.FilledAmount
		item.Type = cell.Type
		item.Amount = cell.Amount
		item.CreatedSec = cell.CreatedSec

		ret = append(ret, item)
	}

	return ret, nil
}

// 查询订单详情
// strOrderID: 订单ID
// return: OrderReturn对象
func (h *Huobi) GetPointOrders() (*proto.Order, error) {
	//orderReturn := OrderReturn{}

	strRequest := fmt.Sprintf("/v1/points/orders")
	jsonPlaceReturn, err := apiKeyGet(make(map[string]string), strRequest, h.accessKey, h.secretKey)
	if nil != err {
		log.Printf("Huobi.GetPointOrders - apiKeyGet failed : %v", err)
		return nil, err
	}
	log.Println(jsonPlaceReturn)
	//json.Unmarshal([]byte(jsonPlaceReturn), &orderReturn)
	//
	//if orderReturn.Status != "ok" {
	//  return nil, errors.New(orderReturn.ErrMsg)
	//}
	//
	//var ret proto.Order
	//ret.Price, _ = strconv.ParseFloat(orderReturn.Data.Price, 64)
	//ret.ID = orderReturn.Data.ID
	//ret.Symbol = orderReturn.Data.Symbol
	//ret.State = orderReturn.Data.State
	//ret.FieldAmount, _ = strconv.ParseFloat(orderReturn.Data.FieldAmount, 64)
	//ret.Type = orderReturn.Data.Type
	//ret.Amount, _ = strconv.ParseFloat(orderReturn.Data.Amount, 64)

	return nil, nil

}

func (h *Huobi) SetDetailListener(listener DetailListener) {
	h.detailListener = listener
}

func (h *Huobi) SetDepthlListener(listener DepthlListener) {
	h.depthListener = listener
}

func (h *Huobi) SetKLineUpListener(listener KLineUpListener) {
	h.klineUpListener = listener
}

func (h *Huobi) SetOrderListener(listener OrderListener) {
	h.orderListener = listener
}

// Listener 订阅事件监听器
type DetailListener func(symbol string, detail *MarketTradeDetail)

func (h *Huobi) SubscribeDetail(symbols ...string) (err error) {
	for _, symbol := range symbols {
		err = h.market.Subscribe("market."+symbol+".trade.detail", func(topic string, j *simplejson.Json) {
			js, _ := j.MarshalJSON()
			var mtd MarketTradeDetail
			err := json.Unmarshal(js, &mtd)
			if err != nil {
				log.Printf("Huobi.SubscribeDetail - callback failed : %s", err.Error())
				return
			}

			ts := strings.Split(topic, ".")
			if h.detailListener != nil {
				h.detailListener(ts[1], &mtd)
			}
		})
		if nil != err {
			log.Printf("Huobi.SubscribeDetail - Subscribe failed : %s", err.Error())
			return err
		}
	}
	return nil
}

// Listener 订阅事件监听器
type DepthlListener func(symbol string, depth *MarketDepth)

func (h *Huobi) SubscribeDepth(symbols ...string) {
	for _, symbol := range symbols {
		h.market.Subscribe("market."+symbol+".depth.step0", func(topic string, j *simplejson.Json) {
			js, _ := j.MarshalJSON()
			var md = MarketDepth{}
			err := json.Unmarshal(js, &md)
			if err != nil {
				log.Printf("Huobi.SubscribeDepth - callback failed : %s", err.Error())
				return
			}

			ts := strings.Split(topic, ".")
			if h.depthListener != nil {
				h.depthListener(ts[1], &md)
			}
		})
	}
}

// KLineUpListener :
type KLineUpListener func(symbol string, kline *KLineUpdate)

// SubscribeKLine :
func (h *Huobi) SubscribeKLine(period string, symbols ...string) {
	for _, symbol := range symbols {
		h.market.Subscribe("market."+symbol+".kline."+period, func(topic string, j *simplejson.Json) {
			js, _ := j.MarshalJSON()
			var mku KLineUpdate
			err := json.Unmarshal(js, &mku)
			if err != nil {
				log.Printf("Huobi.SubscribeKLine - callback failed : %s", err.Error())
				return
			}

			ts := strings.Split(topic, ".")
			if h.klineUpListener != nil {
				h.klineUpListener(ts[1], &mku)
			}
		})
	}
	return
}

// OrderUpdateData :
type OrderUpdateData struct {
	AccountID        int    `json:"account-id"`
	CreatedAt        int    `json:"created-at"`
	FilledAmount     string `json:"filled-amount"`
	FilledCashAmount string `json:"filled-cash-amount"`
	FilledFees       string `json:"filled-fees"`
	OrderAmount      string `json:"order-amount"`
	OrderID          int    `json:"order-id"`
	OrderPrice       string `json:"order-price"`
	OrderSource      string `json:"order-source"`
	OrderState       string `json:"order-state"`
	OrderType        string `json:"order-type"`
	Price            string `json:"price"`
	Role             string `json:"role"`
	SeqID            int    `json:"seq-id"`
	Symbol           string `json:"symbol"`
	UnfilledAmount   string `json:"unfilled-amount"`
}

// OrderUpdate :
type OrderUpdate struct {
	OP    string          `json:"op"`
	Topic string          `json:"topic"`
	ts    int64           `json:"ts"`
	Order OrderUpdateData `json:"data"`
}

// OrderListener : 订阅事件监听器
type OrderListener func(symbol string, od *OrderUpdate)

// SubscribeOrder :
func (h *Huobi) SubscribeOrder(symbols ...string) (err error) {
	for _, symbol := range symbols {
		tp := "order." + symbol
		err = h.market.SubscribeEx(tp,
			map[string]string{"op": "sub", "cid": tp, "topic": tp},
			func(topic string, j *simplejson.Json) {
				js, _ := j.MarshalJSON()
				var order OrderUpdate
				err := json.Unmarshal(js, &order)
				if err != nil {
					log.Printf("Huobi.SubscribeDetail - callback failed : %s", err.Error())
					return
				}

				ts := strings.Split(topic, ".")
				if h.detailListener != nil {
					h.orderListener(ts[1], &order)
				}
			})
		if nil != err {
			log.Printf("Huobi.SubscribeOrder - Subscribe failed : %s", err.Error())
			return err
		}
	}
	return nil
}

// GetMarginLoanOrders :
// 参数名称 是否必须    类型    描述    默认值  取值范围
// symbol   true        string  交易对
// start-date false     string  查询开始日期, 日期格式yyyy-mm-dd
// end-date false       string  查询结束日期, 日期格式yyyy-mm-dd
// states   false       string  状态
// from     false       string  查询起始 ID
// direct   false       string  查询方向        prev 向前，next 向后
// size     false       string  查询记录大小
// 只需要"symbol":"xxxxxx"和 "states":"accrual"， 其他应该都不需要
func (h *Huobi) GetMarginLoanOrders(params map[string]string) ([]LoanOrder, error) {
	strReqURL := "/v1/margin/loan-orders"
	buf, err := apiKeyGet(params, strReqURL, h.accessKey, h.secretKey)
	if nil != err {
		log.Printf("Huobi.GetMarginLoanOrders - apiKeyGet failed : %v", err)
		return nil, err
	}
	var resp LoanOrderResp
	err = json.Unmarshal([]byte(buf), &resp)
	if nil != err {
		log.Printf("Huobi.GetMarginLoanOrders - json.Unmarshal '%s' failed : %v", buf, err)
		return nil, err
	}

	if resp.Status != "ok" {
		log.Printf("Huobi.GetMarginLoanOrders - resp.Status %s, invalid", resp.Status)
		return nil, fmt.Errorf("status %s invalid", resp.Status)
	}

	return resp.Orders, nil
}

// MarginIO :
func (h *Huobi) MarginIO(symbol, currency string, dirIn bool, amount float64) error {
	strReqURL := ""
	if dirIn {
		strReqURL = "/v1/dw/transfer-in/margin"
	} else {
		strReqURL = "/v1/dw/transfer-out/margin"
	}

	params := map[string]string{"symbol": symbol,
		"currency": currency,
		"amount":   fmt.Sprintf("%.8f", amount),
	}

	buf, err := apiKeyPost(params, strReqURL, h.accessKey, h.secretKey)
	if nil != err {
		log.Printf("Huobi.MarginIO - apiKeyPost failed : %v", err)
		return err
	}
	resMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(buf), &resMap)
	if nil != err {
		log.Printf("Huobi.MarginIO - json.Unmarshal '%s' failed : %v", buf, err)
		return err
	}
	if resMap["status"] != "ok" {
		log.Printf("Huobi.MarginIO - status invalid, response : %s", buf)
		return fmt.Errorf("status %s invalid", resMap["status"])
	}
	return nil
}

// RepayLoan :
func (h *Huobi) RepayLoan(loanID int, amount string /*not float*/) (err error) {
	strReqURL := fmt.Sprintf("/v1/margin/orders/%d/repay", loanID)
	buf, err := apiKeyPost(map[string]string{"amount": amount}, strReqURL, h.accessKey, h.secretKey)
	if nil != err {
		log.Printf("Huobi.RepayLoan - apiKeyPost failed : %v", err)
		return err
	}
	resMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(buf), &resMap)
	if nil != err {
		log.Printf("Huobi.RepayLoan - json.Unmarshal '%s' failed : %v", buf, err)
		return err
	}
	if resMap["status"] != "ok" {
		log.Printf("Huobi.RepayLoan - status invalid, response : %s", buf)
		return fmt.Errorf("status %s invalid", resMap["status"])
	}
	return nil
}

// ApplyLoan :
func (h *Huobi) ApplyLoan(symbol, currency string, amount float64) (err error) {
	strReqURL := "/v1/margin/orders"
	buf, err := apiKeyPost(map[string]string{
		"symbol":   symbol,
		"currency": currency,
		"amount":   fmt.Sprintf("%.8f", amount)},
		strReqURL, h.accessKey, h.secretKey)

	if nil != err {
		log.Printf("Huobi.ApplyLoan - apiKeyPost failed : %v", err)
		return err
	}
	resMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(buf), &resMap)
	if nil != err {
		log.Printf("Huobi.ApplyLoan - json.Unmarshal '%s' failed : %v", buf, err)
		return err
	}
	if resMap["status"] != "ok" {
		log.Printf("Huobi.ApplyLoan - status invalid, response : %s", buf)
		return fmt.Errorf("status %s invalid", resMap["status"])
	}
	return nil
}

// SymbolDesc :
type SymbolDesc struct {
	BaseCurrency    string `json:"base-currency"`
	QuoteCurrency   string `json:"quote-currency"`
	PricePrecision  int    `json:"price-precision"`
	AmountPrecision int    `json:"amount-precision"`
	SymbolPartition string `json:"symbol-partition"`
	Symbol          string `json:"symbol"`
}

type getSymbolResp struct {
	Status string       `json:"status"`
	Data   []SymbolDesc `json:"data"`
}

// GetSymbols :
func (h *Huobi) GetSymbols() (descs []SymbolDesc, err error) {
	buf, err := apiKeyGet(map[string]string{}, "/v1/common/symbols", h.accessKey, h.secretKey)
	if nil != err {
		log.Printf("Huobi.GetSymbols - apiKeyGet failed : %v", err)
		return nil, err
	}

	var resp getSymbolResp
	if err = json.Unmarshal([]byte(buf), &resp); nil != err {
		log.Printf("Huobi.GetSymbols - json.Unmarshal '%s' failed : %v", buf, err)
		return nil, err
	}
	if resp.Status != "ok" {
		log.Printf("Huobi.GetSymbols - status invalid, response : %s", buf)
		return nil, fmt.Errorf("status %s invalid", resp.Status)
	}
	return resp.Data, nil
}

type getTickersResp struct {
	Status string `json:"status"`
	Ts     int64  `json:"ts"`
	Data   []struct {
		Open   float64 `json:"open"`
		Close  float64 `json:"close"`
		Low    float64 `json:"low"`
		High   float64 `json:"high"`
		Amount float64 `json:"amount"`
		Count  int     `json:"count"`
		Vol    float64 `json:"vol"`
		Symbol string  `json:"symbol"`
	} `json:"data"`
}

// GetTickers :
func (h *Huobi) GetTickers() (tickerMap map[string]*TickData, err error) {
	buf, err := apiKeyGet(map[string]string{}, "/market/tickers", h.accessKey, h.secretKey)
	if nil != err {
		log.Printf("Huobi.GetTickers - apiKeyGet failed : %v", err)
		return nil, err
	}
	var resp getTickersResp
	if err = json.Unmarshal([]byte(buf), &resp); nil != err {
		log.Printf("Huobi.GetTickers - json.Unmarshal '%s' failed : %v", buf, err)
		return nil, err
	}
	if resp.Status != "ok" {
		log.Printf("Huobi.GetTickers - status invalid, response : %s", buf)
		return nil, fmt.Errorf("status %s invalid", resp.Status)
	}
	tickerMap = make(map[string]*TickData)
	for idx := range resp.Data {
		data := &resp.Data[idx]
		tickerMap[data.Symbol] = &TickData{
			Amount:    data.Amount,
			Direction: "",
			Price:     data.Close,
			TS:        resp.Ts,
		}
	}
	return tickerMap, nil
}

// NewHuobi :
func NewHuobi(accesskey, secretkey string) (*Huobi, error) {
	h := &Huobi{
		accessKey: accesskey,
		secretKey: secretkey,
	}

	if accesskey != "" {
		log.Println("init huobi.")
		ret, err := h.GetAccounts()
		if nil != err {
			return nil, err
		}
		if ret.Status != "ok" {
			return nil, errors.New(ret.ErrMsg)
		}

		for _, account := range ret.Data {
			if account.Type == "spot" {
				log.Println("account id:", account.ID)
				h.tradeAccount.ID = account.ID
				h.tradeAccount.Type = account.Type
				h.tradeAccount.State = account.State
				h.tradeAccount.UserID = account.UserID
				break
			}

		}
	}

	log.Println("init huobi success.")

	return h, nil
}
