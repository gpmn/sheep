package okex

import "github.com/gpmn/sheep/proto"

func TransOrderType(t string) string {
	switch t {
	case OrderPlaceTypeBuy:
		return proto.OrderPlaceTypeBuyLimit
	case OrderPlaceTypeSell:
		return proto.OrderPlaceTypeSellLimit
	case OrderPlaceTypeBuyMarket:
		return proto.OrderPlaceTypeBuyMarket
	case OrderPlaceTypeSellMarket:
		return proto.OrderPlaceTypeSellMarket
	case proto.OrderPlaceTypeBuyLimit:
		return OrderPlaceTypeBuy
	case proto.OrderPlaceTypeSellLimit:
		return OrderPlaceTypeSell
	case proto.OrderPlaceTypeBuyMarket:
		return OrderPlaceTypeBuyMarket
	case proto.OrderPlaceTypeSellMarket:
		return OrderPlaceTypeSellMarket
	default:
		return "类型错误"

	}
}

func TransOrderStateFromStatus(s int) string {
	switch s {
	case OrderStatusCancel:
		return proto.OrderStateCanceled
	case OrderStatusFilled:
		return proto.OrderStateFilled
	case OrderStatusPartialFilled:
		return proto.OrderStatePartialFilled
	case OrderStatusUnsettled:
		return proto.OrderStateSubmitted
	default:
		return "类型错误"

	}
}
