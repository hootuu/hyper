package supplyord

import (
	"context"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/hiorder"
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
	if hlog.IsElapseComponent() {
		defer hlog.ElapseWithCtx(ctx, "OrdDeal.After",
			hlog.F(zap.Uint64("ordID", d.ord.ID),
				zap.String("player", d.ord.Payer.MustToID()),
				zap.Any("src", src),
				zap.Any("target", target)),

			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err)}
				}
				return nil
			},
		)()
	}

	go func() {
		if target == Completed {
			//matter := d.ord.Matter
			//计算金额
		}
	}()
	return nil
}
