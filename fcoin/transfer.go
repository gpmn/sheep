package fcoin

import "github.com/gpmn/sheep/proto"

func TransOrderTypeFromProto(t string) (string, string) {
	switch t {
	case proto.OrderPlaceTypeBuyLimit:
		return OrderTypeLimit, OrderSideBuy
	case proto.OrderPlaceTypeSellLimit:
		return OrderTypeLimit, OrderSideSell
	default:
		return "类型错误", t

	}
}

func TransOrderTypeToProto(t, s string) string {
	if t == OrderTypeLimit {
		if s == OrderSideBuy {
			return proto.OrderPlaceTypeBuyLimit
		} else if s == OrderSideSell {
			return proto.OrderPlaceTypeSellLimit
		}
	}

	return "类型错误" + t + s
}

func TransOrderStateFromStatus(s string) string {
	switch s {
	case OrderStateCanceled:
		return proto.OrderStateCanceled
	case OrderStateFilled:
		return proto.OrderStateFilled
	case OrderStatePartialFilled:
		return proto.OrderStatePartialFilled
	case OrderStateSubmitted:
		return proto.OrderStateSubmitted
	default:
		return "类型错误" + s

	}
}
