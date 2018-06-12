package okex

import (
	"net/url"

	"strings"

	"strconv"

	"log"

	"github.com/leek-box/sheep/consts"
	"github.com/leek-box/sheep/proto"
	"github.com/leizongmin/huobiapi"
	"github.com/pkg/errors"
)

type OKEX struct {
	accessKey string
	secretKey string
	market    *Market
}

func (o *OKEX) GetExchangeType() string {
	return consts.ExchangeTypeOKEX
}

func (o *OKEX) GetAccountBalance() ([]proto.AccountBalance, error) {
	path := "userinfo.do"
	var ret BalanceReturn
	err := o.apiKeyPost(url.Values{}, path, &ret)
	if err != nil {
		return nil, err
	}

	if !ret.Result {
		return nil, errors.New("result error")
	}

	var res []proto.AccountBalance
	for k, v := range ret.Info.Funds.Free {
		if v != "0" {
			var item proto.AccountBalance
			item.Currency = k
			item.Balance = v
			item.Type = proto.AccountBalanceTypeTrade

			res = append(res, item)
		}
	}

	for k, v := range ret.Info.Funds.Freezed {
		if v != "0" {
			var item proto.AccountBalance
			item.Currency = k
			item.Balance = v
			item.Type = proto.AccountBalanceTypeFrozen

			res = append(res, item)
		}
	}

	return res, nil
}

//访问频率 20次/2秒
func (o *OKEX) OrderPlace(params *proto.OrderPlaceParams) (*proto.OrderPlaceReturn, error) {
	path := "trade.do"

	values := url.Values{}
	values.Set("symbol", strings.ToLower(params.BaseCurrencyID)+"_"+strings.ToLower(params.QuoteCurrencyID))
	values.Set("type", TransOrderType(params.Type))
	values.Set("price", strconv.FormatFloat(params.Price, 'f', -1, 64))
	values.Set("amount", strconv.FormatFloat(params.Amount, 'f', -1, 64))

	var okRet OrderPlaceReturn
	err := o.apiKeyPost(values, path, &okRet)
	if err != nil {
		return nil, err
	}
	if okRet.ErrorCode != 0 {
		return nil, codeError(okRet.ErrorCode)
	}
	var ret proto.OrderPlaceReturn

	ret.OrderID = strconv.FormatInt(okRet.OrderID, 10)

	return &ret, nil
}

func (o *OKEX) OrderCancel(params *proto.OrderCancelParams) error {
	path := "cancel_order.do"
	values := url.Values{}
	values.Set("order_id", params.OrderID)
	values.Set("symbol", strings.ToLower(params.BaseCurrencyID)+"_"+strings.ToLower(params.QuoteCurrencyID))

	var okRet CancelOrderReturn
	err := o.apiKeyPost(values, path, &okRet)
	if err != nil {
		return err
	}
	if okRet.ErrorCode != 0 {
		return codeError(okRet.ErrorCode)
	}
	if !okRet.Result {
		return errors.New("撤单失败")
	}

	return nil

}

func (o *OKEX) GetOrderInfo(params *proto.OrderInfoParams) (*proto.Order, error) {
	path := "order_info.do"
	values := url.Values{}
	values.Set("symbol", strings.ToLower(params.BaseCurrencyID)+"_"+strings.ToLower(params.QuoteCurrencyID))
	values.Set("order_id", params.OrderID)

	var okRet OrderInfoReturn
	err := o.apiKeyPost(values, path, &okRet)
	if err != nil {
		return nil, err
	}
	if okRet.ErrorCode != 0 {
		return nil, codeError(okRet.ErrorCode)
	}
	if !okRet.Result {
		return nil, errors.New("获取失败")
	}
	if len(okRet.Orders) == 0 {
		return nil, errors.New("获取失败")
	}

	okOrder := okRet.Orders[0]

	var ret proto.Order
	ret.ID = strconv.FormatInt(okOrder.OrderID, 10)
	ret.Symbol = strings.Replace(okOrder.Symbol, "_", "", 1)
	ret.State = TransOrderStateFromStatus(okOrder.Status)
	ret.Amount = okOrder.Amount
	ret.FieldAmount = okOrder.DealAmount
	ret.Price = okOrder.Price
	ret.Type = TransOrderType(okOrder.Type)

	return &ret, nil
}

func (o *OKEX) GetOrders(params *proto.OrdersParams) ([]proto.Order, error) {
	path := "order_history.do"
	values := url.Values{}
	values.Set("symbol", strings.ToLower(params.BaseCurrencyID)+"_"+strings.ToLower(params.QuoteCurrencyID))
	values.Set("status", params.Status)
	values.Set("current_page", params.CurrentPage)
	values.Set("page_length", params.PageLength)

	var okRet OrderInfoReturn
	err := o.apiKeyPost(values, path, &okRet)
	if err != nil {
		return nil, err
	}
	if okRet.ErrorCode != 0 {
		return nil, codeError(okRet.ErrorCode)
	}
	if !okRet.Result {
		return nil, errors.New("获取失败")
	}

	var ret []proto.Order
	for _, okOrder := range okRet.Orders {
		var item proto.Order
		item.ID = strconv.FormatInt(okOrder.OrderID, 10)
		item.Symbol = strings.Replace(okOrder.Symbol, "_", "", 1)
		item.State = TransOrderStateFromStatus(okOrder.Status)
		item.Amount = okOrder.Amount
		item.FieldAmount = okOrder.DealAmount
		item.Price = okOrder.Price
		item.Type = TransOrderType(okOrder.Type)

		ret = append(ret, item)
	}

	return ret, nil

}

func (h *OKEX) SubscribeDetail(symbols ...string) {
	for _, symbol := range symbols {
		h.market.Subscribe("ok_sub_spot_"+symbol+"_depth", func(topic string, j *huobiapi.JSON) {
			js, _ := j.MarshalJSON()
			log.Println(js)

		})
	}

}

func NewOKEX(apiKey, secretKey string) (*OKEX, error) {
	o := &OKEX{
		accessKey: apiKey,
		secretKey: secretKey,
	}

	var err error
	o.market, err = NewMarket()
	if err != nil {
		return nil, err
	}

	go o.market.Loop()

	return o, nil
}
