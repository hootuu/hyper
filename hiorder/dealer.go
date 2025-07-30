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

type EmptyDealer[T Matter] struct {
	code Code
}

func NewEmptyDealer[T Matter](code Code) *EmptyDealer[T] {
	return &EmptyDealer[T]{code}
}

func (e *EmptyDealer[T]) Code() Code {
	return e.code
}

func (e *EmptyDealer[T]) Build(ord Order[T]) (Deal[T], error) {
	return &EmptyDeal[T]{
		dealer: e,
	}, nil
}

func (e *EmptyDealer[T]) OnPaymentAltered(_ context.Context, _ *payment.AlterPayload) error {
	return nil
}

func (e *EmptyDealer[T]) OnShippingAltered(_ context.Context, _ *shipping.AlterPayload) error {
	return nil
}

type EmptyDeal[T Matter] struct {
	dealer *EmptyDealer[T]
}

func (e *EmptyDeal[T]) Code() Code {
	return e.dealer.Code()
}

func (e *EmptyDeal[T]) Before(_ context.Context, _ Status, _ Status) error {
	return nil
}

func (e *EmptyDeal[T]) After(_ context.Context, _ Status, _ Status) error {
	return nil
}
