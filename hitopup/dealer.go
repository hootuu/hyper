package hitopup

import (
	"context"
	"fmt"
	"github.com/hootuu/hyle/hcoin"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/hiorder"
	"go.uber.org/zap"
	"time"
)

type Dealer struct {
	code     hiorder.Code
	currency hcoin.Currency
	timeout  time.Duration

	f     *hiorder.Factory[Matter]
	topup *TopUp
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

func (d *Dealer) doInit(f *hiorder.Factory[Matter], topup *TopUp) {
	d.topup = topup
	d.f = f
}

func (d *Dealer) Code() hiorder.Code {
	return d.code
}

func (d *Dealer) Currency() hcoin.Currency {
	return d.currency
}

func (d *Dealer) Build(ord hiorder.Order[Matter]) (hiorder.Deal[Matter], error) {
	return newDeal(&ord, d, d.topup), nil
}

// OnPaymentAltered todo use ctx for mq
func (d *Dealer) OnPaymentAltered(alter *hiorder.PaymentAltered[Matter]) error {

	if alter == nil {
		hlog.Fix("hitopup.dealer.OnPaymentAltered: alter is nil")
		return nil
	}
	switch alter.DstStatus {
	case hiorder.PaymentPaid:
		ctx := context.Background()
		eng, err := d.f.Load(ctx, alter.Order.ID)
		if err != nil {
			hlog.Err("hitopup.dealer.OnPaymentAltered: load fail", zap.Error(err))
			return err
		}
		fmt.Println("hitopup.dealer.OnPaymentAltered: load success") // todo
		err = eng.Complete(ctx)
		if err != nil {
			hlog.Err("hitopup.dealer.OnPaymentAltered: complete fail", zap.Error(err))
			return err
		}
		return nil
	default:
		return nil
	}
}
