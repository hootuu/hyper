package hitopup

import (
	"context"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/hiorder"
	"github.com/nineora/harmonic/harmonic"
	"github.com/nineora/harmonic/nineapi"
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
	case hiorder.Completed:
		nine := harmonic.Nineora()
		sig, err := nine.TokenMint(ctx, &nineapi.TxMintParas{
			Mint:       d.topup.mint,
			Recipient:  d.ord.Matter.InAccount,
			Amount:     d.ord.Amount,
			LockAmount: 0,
			//Meta:       d.ord.Meta, TODO
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
