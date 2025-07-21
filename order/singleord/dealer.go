package singleord

import (
	"github.com/hootuu/hyper/hiorder"
	"time"
)

type Dealer struct {
	code    hiorder.Code
	timeout time.Duration

	f      *hiorder.Factory[Matter]
	single *Single
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

func (d *Dealer) Code() hiorder.Code {
	return d.code
}

func (d *Dealer) Build(ord hiorder.Order[Matter]) (hiorder.Deal[Matter], error) {
	return newDeal(d, &ord), nil
}

func (d *Dealer) OnPaymentAltered(alter *hiorder.PaymentAltered[Matter]) error {
	return nil
}

func (d *Dealer) doInit(f *hiorder.Factory[Matter], single *Single) {
	d.single = single
	d.f = f
}
