package fcoin

import (
	"encoding/json"
	"fmt"
	"github.com/leek-box/sheep/proto"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

const FCoinHost = "https://api.fcoin.com/v2/"

type FCoin struct {
	accessKey string
	secretKey string
	Market    *Market
}

func (f *FCoin) OpenWebsocket() error {
	var err error
	f.Market, err = NewMarket()
	if err != nil {
		return err
	}

	go f.Market.Loop()
	return nil
}

func (f *FCoin) GetAccountBalance() ([]proto.AccountBalance, error) {
	balanceReturn := BalanceReturn{}
	strRequest := "accounts/balance"
	jsonBanlanceReturn := apiKeyGet(make(map[string]string), strRequest, f.accessKey, f.secretKey)
	json.Unmarshal([]byte(jsonBanlanceReturn), &balanceReturn)
	if balanceReturn.Status != 0 {
		return nil, errors.New(strconv.Itoa(balanceReturn.Status))
	}

	var res []proto.AccountBalance
	for _, blance := range balanceReturn.Data {
		var item proto.AccountBalance
		item.Currency = blance.Currency
		item.Balance = blance.Balance
		item.Type = "trade"

		res = append(res, item)

		var item2 proto.AccountBalance
		item2.Currency = blance.Currency
		item2.Balance = blance.Frozen
		item2.Type = "frozen"

		res = append(res, item2)
	}

	return res, nil
}

// 下单
// placeRequestParams: 下单信息
// return: OrderID
func (f *FCoin) OrderPlace(params *proto.OrderPlaceParams) (*proto.OrderPlaceReturn, error) {
	placeReturn := PlaceReturn{}
	var placeRequestParams PlaceRequestParams
	placeRequestParams.Amount = strconv.FormatFloat(params.Amount, 'f', -1, 64)
	placeRequestParams.Price = strconv.FormatFloat(params.Price, 'f', -1, 64)
	placeRequestParams.Symbol = strings.ToLower(params.BaseCurrencyID) + strings.ToLower(params.QuoteCurrencyID)
	placeRequestParams.Type, placeRequestParams.Side = TransOrderTypeFromProto(params.Type)

	mapParams := make(map[string]string)
	mapParams["amount"] = placeRequestParams.Amount
	if 0 < len(placeRequestParams.Price) {
		mapParams["price"] = placeRequestParams.Price
	}
	mapParams["symbol"] = placeRequestParams.Symbol
	mapParams["type"] = placeRequestParams.Type
	mapParams["side"] = placeRequestParams.Side

	strRequest := "orders"
	jsonPlaceReturn := apiKeyPost(mapParams, strRequest, f.accessKey, f.secretKey)
	json.Unmarshal([]byte(jsonPlaceReturn), &placeReturn)

	if placeReturn.Status != 0 {
		return nil, errors.New(strconv.Itoa(placeReturn.Status))
	}

	var ret proto.OrderPlaceReturn
	ret.OrderID = placeReturn.Data

	return &ret, nil

}

// 申请撤销一个订单请求
// strOrderID: 订单ID
// return: PlaceReturn对象
func (f *FCoin) OrderCancel(params *proto.OrderCancelParams) error {
	placeReturn := PlaceReturn{}

	strRequest := fmt.Sprintf("orders/%s/submit-cancel", params.OrderID)
	jsonPlaceReturn := apiKeyPost(make(map[string]string), strRequest, f.accessKey, f.secretKey)
	json.Unmarshal([]byte(jsonPlaceReturn), &placeReturn)

	return nil
}

// 查询订单详情
// strOrderID: 订单ID
// return: OrderReturn对象
func (f *FCoin) GetOrderInfo(params *proto.OrderInfoParams) (*proto.Order, error) {
	orderReturn := OrderReturn{}

	strRequest := fmt.Sprintf("orders/%s", params.OrderID)
	jsonPlaceReturn := apiKeyGet(make(map[string]string), strRequest, f.accessKey, f.secretKey)
	json.Unmarshal([]byte(jsonPlaceReturn), &orderReturn)

	var ret proto.Order
	ret.Price, _ = strconv.ParseFloat(orderReturn.Data.Price, 64)
	ret.ID = orderReturn.Data.ID
	ret.Symbol = orderReturn.Data.Symbol
	ret.State = TransOrderStateFromStatus(orderReturn.Data.State)
	ret.Type = TransOrderTypeToProto(orderReturn.Data.Type, orderReturn.Data.Side)
	ret.Amount, _ = strconv.ParseFloat(orderReturn.Data.Amount, 64)

	return &ret, nil

}

func (f *FCoin) GetOrders(params *proto.OrdersParams) ([]proto.Order, error) {
	ordersReturn := OrdersReturn{}

	var paramMap = make(map[string]string)
	paramMap["symbol"] = params.Symbol
	paramMap["states"] = params.States
	//paramMap["before"] = ""
	//paramMap["after"] = ""
	paramMap["limit"] = "10"

	strRequest := "orders"
	jsonRet := apiKeyGet(paramMap, strRequest, f.accessKey, f.secretKey)
	json.Unmarshal([]byte(jsonRet), &ordersReturn)

	var ret []proto.Order
	for _, cell := range ordersReturn.Data {
		var item proto.Order
		item.Price, _ = strconv.ParseFloat(cell.Price, 64)
		item.ID = cell.ID
		item.Symbol = cell.Symbol
		item.State = cell.State
		item.FieldAmount, _ = strconv.ParseFloat(cell.FilledAmount, 64)
		item.Type = TransOrderTypeToProto(cell.Type, cell.Side)
		item.Amount, _ = strconv.ParseFloat(cell.Amount, 64)

		ret = append(ret, item)
	}

	return ret, nil

}

func GetMarketDepth(params *proto.MarketDepthParams) (*MarketDepthReturn, error) {
	marketDepth := MarketDepthReturn{}

	strRequest := "market/depth/" + params.Level + "/" + params.Symbol
	jsonRet := apiKeyGet(make(map[string]string), strRequest, "", "")

	json.Unmarshal([]byte(jsonRet), &marketDepth)

	return &marketDepth, nil
}

func (f *FCoin) CloseWebsocket() error {
	return f.Market.Close()
}

func NewFCoin(accessKey, secretKey string) (*FCoin, error) {
	if accessKey == "" || secretKey == "" {
		return nil, errors.New("access key or secret key error")
	}
	f := &FCoin{
		accessKey: accessKey,
		secretKey: secretKey,
	}

	return f, nil
}
