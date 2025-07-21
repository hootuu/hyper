package hiorder

import (
	"context"
	"time"
)

type Dealer[T Matter] interface {
	Code() Code
	Build(ord Order[T]) (Deal[T], error)
	OnPaymentAltered(alter *PaymentAltered[T]) error
}

type Deal[T Matter] interface {
	Code() Code
	Timeout() time.Duration
	Before(ctx context.Context, src Status, target Status) error
	After(ctx context.Context, src Status, target Status) error
}
