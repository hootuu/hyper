package hitopup

import (
	"context"
	"github.com/hootuu/hyper/hiorder"
	"github.com/hootuu/hyper/payment"
	"github.com/hootuu/hyper/shipping"
	"time"
)

type Dealer struct {
	code    hiorder.Code
	timeout time.Duration

	f     *hiorder.Factory[Matter]
	topup *TopUp
}

func newDealer(
	code hiorder.Code,
	timeout time.Duration,
) *Dealer {
	return &Dealer{
		code:    code,
		timeout: timeout,
	}
}

func (d *Dealer) doInit(f *hiorder.Factory[Matter], topup *TopUp) {
	d.topup = topup
	d.f = f
}

func (d *Dealer) Code() hiorder.Code {
	return d.code
}

func (d *Dealer) Build(ord hiorder.Order[Matter]) (hiorder.Deal[Matter], error) {
	return newDeal(&ord, d, d.topup), nil
}

func (d *Dealer) OnPaymentAltered(_ context.Context, _ *payment.AlterPayload) error {
	return nil
}

func (d *Dealer) OnShippingAltered(_ context.Context, _ *shipping.AlterPayload) error {
	return nil
}
