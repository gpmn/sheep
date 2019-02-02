package sheep

import (
	"github.com/gpmn/sheep/consts"
	"github.com/gpmn/sheep/huobi"
	"github.com/gpmn/sheep/okex"
	"github.com/gpmn/sheep/proto"
	"github.com/pkg/errors"
)

type ExchageI interface {
	GetExchangeType() string
	//获取账户余额
	GetAccountBalance() ([]proto.AccountBalance, error)
	//下单
	OrderPlace(params *proto.OrderPlaceParams) (*proto.OrderPlaceReturn, error)
	//取消订单
	OrderCancel(params *proto.OrderCancelParams) error
	//获取订单详情
	GetOrderInfo(params *proto.OrderInfoParams) (*proto.Order, error)
	//获取历史订单列表
	GetOrders(params *proto.OrdersParams) ([]proto.Order, error)
}

func NewExchange(typ, accessKey, secretKey string) (ExchageI, error) {
	switch typ {
	case consts.ExchangeTypeHuobi:
		return huobi.NewHuobi(accessKey, secretKey)
	case consts.ExchangeTypeOKEX:
		return okex.NewOKEX(accessKey, secretKey)
	}

	return nil, errors.New("不支持该交易所")
}
