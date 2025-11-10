package gjjord

import (
	"context"
	"github.com/hootuu/hyper/hiorder"
	"github.com/hootuu/hyper/payment"
	"github.com/hootuu/hyper/shipping"
)

type Dealer struct {
	code hiorder.Code
	f    *Factory
}

func newDealer(code hiorder.Code, f *Factory) *Dealer {
	return &Dealer{
		code: code,
		f:    f,
	}
}

func (d *Dealer) Code() hiorder.Code {
	return d.code
}

func (d *Dealer) Build(ord hiorder.Order[Matter]) (hiorder.Deal[Matter], error) {
	return newDeal(d, &ord), nil
}

func (d *Dealer) OnPaymentAltered(ctx context.Context, alter *payment.AlterPayload) (err error) {
	return nil
}

func (d *Dealer) OnShippingAltered(ctx context.Context, alter *shipping.AlterPayload) (err error) {
	return nil
}
