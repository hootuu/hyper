package gjjord

import (
	"context"
	"fmt"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/ex"
	"github.com/hootuu/hyper/hiorder"
	"github.com/nineora/lightv/lightv"
	"github.com/nineora/lightv/qing"
	"github.com/spf13/cast"
	"go.uber.org/zap"
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

	if d.ord.Ex.Meta.Get("is_promotion").Bool() {
		eng, err := GetFactory().Core().Load(ctx, d.ord.ID)
		if err != nil {
			return err
		}
		if eng.GetOrder().Status == hiorder.Consensus {
			return eng.AdvToCompleted(ctx)
		}
		return nil
	} else {
		go func(d *Deal, ctx context.Context, target hiorder.Status) {
			var awErr error
			defer hlog.EL(ctx, "gjjord.award.After goroutine").
				With(
					zap.Uint64("orderID", d.ord.ID),
					zap.String("buyer", d.ord.Payer.MustToID()),
					zap.String("targetStatus", cast.ToString(target)),
					zap.Uint64("orderAmount", d.ord.Amount),
					zap.String("biz", d.Code()),
				).
				EndWith(func() []zap.Field {
					return []zap.Field{
						zap.Error(awErr),
					}
				})()

			orderID := d.ord.ID
			buyer := d.ord.Payer.MustToID()
			amount := d.ord.Amount

			newEx := ex.NewEx()
			newEx.Meta.Set("player", buyer).Set("order_id", orderID)

			orderM, err := hiorder.DbMustGet(ctx, cast.ToString(orderID))
			if err != nil {
				hlog.TraceErr("gjjord.Deal.After: DbMustGet failed", ctx, err)
				return
			}
			if orderM.ConsensusTime == nil {
				return
			}
			// 2025-12-27 17:00:00 之前的订单，奖励的道具都需要冻结和解冻
			limitTime := time.Date(2025, 12, 30, 11, 0, 0, 0, time.Local)
			isLock := orderM.ConsensusTime.Before(limitTime)

			if target == hiorder.Consensus {
				cost := cast.ToUint64(d.ord.Ex.Meta.Get("product.cost").Data())
				totalCost := cost * d.ord.Matter.Count
				if awErr = lightv.AwardOrderPrepare(ctx, isLock, orderID, buyer, amount-totalCost, amount, d.Code(), newEx); awErr != nil {
					hlog.TraceErr("gjjord.Deal.After: AwardOrderPrepare failed", ctx, awErr)
					return
				}
			}

			if target == hiorder.Refunded {
				if awErr = lightv.AwardOrderCancel(ctx, isLock, orderID, buyer, d.Code(), newEx); awErr != nil {
					hlog.TraceErr("gjjord.Deal.After: AwardOrderCancel failed", ctx, awErr)
					return
				}

			}
			if target == hiorder.Completed {
				if awErr = lightv.AwardOrderConfirm(ctx, isLock, orderID, buyer, d.Code(), newEx); awErr != nil {
					hlog.TraceErr("gjjord.Deal.After: AwardOrderComplete failed", ctx, awErr)
					return
				}
				if awErr = lightv.Assets.TerrTaxing(ctx, cast.ToString(orderID), d.Code(), buyer, amount, newEx, 10000, qing.GJJTerrRadioMap); awErr != nil {
					hlog.TraceFix(fmt.Sprintf("lightv.AwardByOrder: TerrTaxing failed for order %d", orderID), ctx, awErr, zap.Uint64("orderID", orderID))
					return
				}
			}
		}(d, ctx, target)
	}
	return nil
}
