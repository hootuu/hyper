package gjjord

import (
	"context"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/hiorder"
	"github.com/nineora/lightv/lightv"
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

func (d *Deal) Before(ctx context.Context, src hiorder.Status, target hiorder.Status) error {
	return nil
}

func (d *Deal) After(ctx context.Context, src hiorder.Status, target hiorder.Status) (err error) {
	if target == hiorder.Consensus {
		go func(d *Deal) {
			if err := lightv.Assets.AwardByOrder(ctx, d.ord.ID, d.ord.Payer.MustToID(), d.ord.Amount, d.Code(), nil); err != nil {
				hlog.TraceErr("gjjord.Deal.After: AwardByOrder failed", ctx, err)
			}
		}(d)
		return nil
	}
	return nil
}
