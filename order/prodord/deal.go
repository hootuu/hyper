package prodord

import (
	"context"
	"fmt"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/hiorder"
	"github.com/hootuu/hyper/hiprod"
	ordtypes "github.com/hootuu/hyper/order/types"
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
		matter := d.ord.Matter

		if target == hiorder.Timeout {
			for _, item := range matter.Items {
				if item.VwhID == 0 || item.SkuID == 0 || item.Quantity == 0 {
					hlog.TraceFix("OrdDeal.After: vwhID or skuID or quantity is zero, skip inventory return", ctx, nil)
					continue
				}
				err = hiprod.SkuStockReset(ctx, hiprod.SkuStockResetParas{
					Vwh:      item.VwhID,
					Pwh:      item.PwhID,
					Sku:      item.SkuID,
					Quantity: item.Quantity,
				})
				if err != nil {
					hlog.TraceFix(fmt.Sprintf("OrdDeal.After: sku stock reset failed, vwhID: %d, skuID: %d", item.VwhID, item.SkuID), ctx, nil)
				}
			}
		} else if target == hiorder.Consensus || target == hiorder.Executing || target == hiorder.Completed {
			mqPublishOrderAlter(&ordtypes.AlterPayload{
				OrderID: d.ord.ID,
				Code:    d.dealer.code,
				Src:     src,
				Dst:     target,
			})
		}
	}()
	return nil
}
