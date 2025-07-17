package singleord

import (
	"context"
	"github.com/hootuu/hyle/hcoin"
	"github.com/hootuu/hyper/hiorder"
	"time"
)

type Deal struct {
	dealer *Dealer
	ord    *hiorder.Order[Matter]
}

func newDeal(dealer *Dealer, ord *hiorder.Order[Matter]) *Deal {
	return &Deal{
		dealer: dealer,
		ord:    ord,
	}
}

func (d *Deal) Code() hiorder.Code {
	return d.dealer.code
}

func (d *Deal) Currency() hcoin.Currency {
	return d.dealer.currency
}

func (d *Deal) Timeout() time.Duration {
	return d.dealer.timeout
}

func (d *Deal) Before(ctx context.Context, src hiorder.Status, target hiorder.Status) error {
	return nil
}

func (d *Deal) After(ctx context.Context, src hiorder.Status, target hiorder.Status) error {
	return nil
}
