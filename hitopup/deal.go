package hitopup

import (
	"context"
	"fmt"
	"github.com/hootuu/hyle/data/ctrl"
	"github.com/hootuu/hyle/data/tag"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/ex"
	"github.com/hootuu/hyper/hiorder"
	"github.com/nineora/harmonic/harmonic"
	"github.com/nineora/harmonic/nineapi"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"time"
)

type Deal struct {
	dealer *Dealer
	ord    *hiorder.Order[Matter]
	topup  *TopUp
}

func newDeal(ord *hiorder.Order[Matter], dealer *Dealer, topup *TopUp) *Deal {
	return &Deal{
		ord:    ord,
		dealer: dealer,
		topup:  topup,
	}
}

func (d *Deal) Code() hiorder.Code {
	return d.dealer.Code()
}

func (d *Deal) Timeout() time.Duration {
	return d.dealer.timeout
}

func (d *Deal) Before(_ context.Context, _ hiorder.Status, _ hiorder.Status) error {
	return nil
}

func (d *Deal) After(ctx context.Context, _ hiorder.Status, target hiorder.Status) error {
	switch target {
	case hiorder.Consensus:
		eng, err := d.topup.factory.Load(ctx, d.ord.ID)
		if err != nil {
			return err
		}
		return eng.AdvToCompleted(ctx)
	case hiorder.Completed:
		nine := harmonic.Nineora()
		sig, err := nine.TokenMint(ctx, &nineapi.TxMintParas{
			Idem:       fmt.Sprintf("%s:%d", d.Code(), d.ord.ID),
			Mint:       d.topup.mint,
			Recipient:  d.ord.Matter.InAccount,
			Amount:     d.ord.Amount,
			LockAmount: 0,
			Ex: &ex.Ex{
				Ctrl: ctrl.MustNewCtrl().MustSet(8, true),
				Tag:  tag.NewTag(d.ord.Code, cast.ToString(d.ord.Code)),
				Meta: d.ord.GetDigest(),
			},
			Biz:  d.ord.Code,
			Link: d.ord.BuildCollar().Link(),
		})
		if err != nil {
			hlog.Err("hitopup.deal.Alter", zap.Error(err))
			return err
		}
		hlog.Logger().Info("hitopup.deal.Alter: OK", zap.String("sig", string(sig)))
	default:
	}
	return nil
}
