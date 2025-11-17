package gjjord

import (
	"context"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/hiorder"
	"github.com/nineora/lightv/lightv"
	"github.com/spf13/cast"
	"go.uber.org/zap"
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
	if d.ord == nil {
		hlog.TraceFix("gjjord.Deal.After: order is nil", ctx, nil)
		return nil
	}
	defer hlog.EL(ctx, "gjjord.Deal.After").
		With(
			zap.Uint64("orderID", d.ord.ID),
			zap.String("buyer", d.ord.Payer.MustToID()),
			zap.Uint64("orderAmount", d.ord.Amount),
			zap.String("biz", d.Code()),
		).
		EndWith(func() []zap.Field {
			return []zap.Field{
				zap.Error(err),
			}
		})()
	if target == hiorder.Consensus {
		go func(d *Deal) {
			cost := cast.ToUint64(d.ord.Ex.Meta.Get("product.cost").Data())
			totalCost := cost * d.ord.Matter.Count
			if err := lightv.Assets.AwardByOrder(ctx, d.ord.ID, d.ord.Payer.MustToID(), d.ord.Amount-totalCost, d.Code(), nil); err != nil {
				hlog.TraceErr("gjjord.Deal.After: AwardByOrder failed", ctx, err)
			}
		}(d)
		return nil
	}
	return nil
}
