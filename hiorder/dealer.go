package hiorder

import (
	"context"
	"github.com/hootuu/hyle/hcoin"
	"time"
)

type Dealer[T Matter] interface {
	Code() Code
	Currency() hcoin.Currency
	Build(ord Order[T]) (Deal[T], error)
	OnPaymentAltered(alter *PaymentAltered[T]) error
}

type Deal[T Matter] interface {
	Code() Code
	Currency() hcoin.Currency
	Timeout() time.Duration
	Before(ctx context.Context, src Status, target Status) error
	After(ctx context.Context, src Status, target Status) error
}
