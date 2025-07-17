package singleord

import (
	"github.com/hootuu/hyle/hcoin"
	"github.com/hootuu/hyper/hiorder"
	"time"
)

type Dealer struct {
	code     hiorder.Code
	currency hcoin.Currency
	timeout  time.Duration

	f      *hiorder.Factory[Matter]
	single *Single
}

func newDealer(
	code hiorder.Code,
	currency hcoin.Currency,
	timeout time.Duration,
) *Dealer {
	return &Dealer{
		code:     code,
		currency: currency,
		timeout:  timeout,
	}
}

func (d *Dealer) Code() hiorder.Code {
	return d.code
}

func (d *Dealer) Currency() hcoin.Currency {
	return d.currency
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
