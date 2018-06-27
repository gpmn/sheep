package bibox

import (
	"errors"
	"log"

	"encoding/json"

	"github.com/leek-box/sheep/util"
)

const BiboxHost = "https://api.bibox.com"

type Bibox struct {
	accessKey string
	secretKey string
}

func (b *Bibox) GetAccountBalabce() (*GetAccountBalanceRsp, error) {
	path := "/v1/transfer"
	var cmd Cmd
	cmd.Cmd = "transfer/assets"
	cmd.Body = map[string]string{
		"select": "1",
	}

	var cmds []Cmd
	cmds = append(cmds, cmd)
	mcmds, _ := json.Marshal(cmds)

	var req = map[string]string{
		"cmds":   string(mcmds),
		"apikey": b.accessKey,
		"sign":   CreateSign(b.secretKey, string(mcmds)),
	}

	ret := util.HttpPostRequest(BiboxHost+path, req, nil)
	var rsp GetAccountBalanceRsp

	err := json.Unmarshal([]byte(ret), &rsp)
	if err != nil {
		return nil, errors.New(ret)
	}

	log.Println(rsp)

	return &rsp, nil

}

func (b *Bibox) OrderPlace(pair, account_type, order_type, order_side, price, amount string) (*OrderPlaceRsp, error) {
	path := "/v1/orderpending"
	var cmd Cmd
	cmd.Cmd = "orderpending/trade"
	cmd.Body = map[string]string{
		"pair":         pair,
		"account_type": account_type,
		"order_type":   order_type,
		"order_side":   order_side,
		"price":        price,
		"amount":       amount,
	}

	var cmds []Cmd
	cmds = append(cmds, cmd)
	mcmds, _ := json.Marshal(cmds)

	var req = map[string]string{
		"cmds":   string(mcmds),
		"apikey": b.accessKey,
		"sign":   CreateSign(b.secretKey, string(mcmds)),
	}

	ret := util.HttpPostRequest(BiboxHost+path, req, nil)
	log.Println(ret)
	var rsp OrderPlaceRsp

	err := json.Unmarshal([]byte(ret), &rsp)
	if err != nil {
		return nil, errors.New(ret)
	}

	log.Println(rsp)

	return &rsp, nil
}

func (b *Bibox) OrderCancel(orders_id string) error {
	path := "/v1/orderpending"
	var cmd Cmd
	cmd.Cmd = "orderpending/cancelTrade"
	cmd.Body = map[string]string{
		"orders_id": orders_id,
	}

	var cmds []Cmd
	cmds = append(cmds, cmd)
	mcmds, _ := json.Marshal(cmds)

	var req = map[string]string{
		"cmds":   string(mcmds),
		"apikey": b.accessKey,
		"sign":   CreateSign(b.secretKey, string(mcmds)),
	}

	ret := util.HttpPostRequest(BiboxHost+path, req, nil)
	log.Println(ret)
	var rsp OrderCancelRsp

	err := json.Unmarshal([]byte(ret), &rsp)
	if err != nil {
		return errors.New(ret)
	}

	log.Println(rsp)

	return nil
}

func (b *Bibox) GetOrderPendingList(pair, account_type, page, size, coin_symbol, currency_symbol, order_side string) (*OrderPendingListRsp, error) {
	path := "/v1/orderpending"
	var cmd Cmd
	cmd.Cmd = "orderpending/orderPendingList"
	cmd.Body = map[string]string{
		"pair":            pair,
		"account_type":    account_type,
		"page":            page,
		"size":            size,
		"coin_symbol":     coin_symbol,
		"currency_symbol": currency_symbol,
		"order_side":      order_side,
	}

	var cmds []Cmd
	cmds = append(cmds, cmd)
	mcmds, _ := json.Marshal(cmds)

	var req = map[string]string{
		"cmds":   string(mcmds),
		"apikey": b.accessKey,
		"sign":   CreateSign(b.secretKey, string(mcmds)),
	}

	ret := util.HttpPostRequest(BiboxHost+path, req, nil)
	log.Println(ret)
	var rsp OrderPendingListRsp

	err := json.Unmarshal([]byte(ret), &rsp)
	if err != nil {
		return nil, errors.New(ret)
	}

	log.Println(rsp)

	return &rsp, nil
}

func (b *Bibox) GetOrderInfo(id string) (*OrderInfoRsp, error) {
	path := "/v1/orderpending"
	var cmd Cmd
	cmd.Cmd = "orderpending/order"
	cmd.Body = map[string]string{
		"id": id,
	}

	var cmds []Cmd
	cmds = append(cmds, cmd)
	mcmds, _ := json.Marshal(cmds)

	var req = map[string]string{
		"cmds":   string(mcmds),
		"apikey": b.accessKey,
		"sign":   CreateSign(b.secretKey, string(mcmds)),
	}

	ret := util.HttpPostRequest(BiboxHost+path, req, nil)
	log.Println(ret)
	var rsp OrderInfoRsp

	err := json.Unmarshal([]byte(ret), &rsp)
	if err != nil {
		return nil, errors.New(ret)
	}

	log.Println(rsp)

	return &rsp, nil
}

func (b *Bibox) GetOrderHistoryList(pair, account_type, page, size, coin_symbol, currency_symbol, order_side string) (*GetOrderHistoryListRsp, error) {
	path := "/v1/orderpending"
	var cmd Cmd
	cmd.Cmd = "orderpending/orderHistoryList"
	cmd.Body = map[string]string{
		"pair":            pair,
		"account_type":    account_type,
		"page":            page,
		"size":            size,
		"coin_symbol":     coin_symbol,
		"currency_symbol": currency_symbol,
		"order_side":      order_side,
	}

	var cmds []Cmd
	cmds = append(cmds, cmd)
	mcmds, _ := json.Marshal(cmds)

	var req = map[string]string{
		"cmds":   string(mcmds),
		"apikey": b.accessKey,
		"sign":   CreateSign(b.secretKey, string(mcmds)),
	}

	ret := util.HttpPostRequest(BiboxHost+path, req, nil)
	log.Println(ret)
	var rsp GetOrderHistoryListRsp

	err := json.Unmarshal([]byte(ret), &rsp)
	if err != nil {
		return nil, errors.New(ret)
	}

	log.Println(rsp)

	return &rsp, nil
}

func NewBibox(accessKey, secretKey string) (*Bibox, error) {
	if accessKey == "" || secretKey == "" {
		return nil, errors.New("access key or secret key error")
	}
	f := &Bibox{
		accessKey: accessKey,
		secretKey: secretKey,
	}

	return f, nil
}

func Ping() {
	path := "/v1/public"
	var req = map[string]string{
		"cmd": "ping",
	}

	ret := util.HttpGetRequest(BiboxHost+path+"?"+util.Map2UrlQuery(req), nil)

	log.Println(ret)

}

func GetMarketDepth(pair string) (*GetMarketDepthRsp, error) {
	path := "/v1/mdata"
	var req = map[string]string{
		"cmd":  "depth",
		"pair": pair,
		"size": "10",
	}

	ret := util.HttpGetRequest(BiboxHost+path+"?"+util.Map2UrlQuery(req), nil)

	var rsp GetMarketDepthRsp

	err := json.Unmarshal([]byte(ret), &rsp)
	if err != nil {
		return nil, err
	}

	return &rsp, nil
}
