package hiorder

import (
	"context"
	"github.com/hootuu/hyper/payment"
	"github.com/hootuu/hyper/shipping"
)

type Dealer[T Matter] interface {
	Code() Code
	Build(ord Order[T]) (Deal[T], error)
	OnPaymentAltered(ctx context.Context, payload *payment.AlterPayload) error
	OnShippingAltered(ctx context.Context, payload *shipping.AlterPayload) error
}

type Deal[T Matter] interface {
	Code() Code
	Before(ctx context.Context, src Status, target Status) error
	After(ctx context.Context, src Status, target Status) error
}
