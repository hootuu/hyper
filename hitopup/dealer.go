package hitopup

import (
	"context"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/hiorder"
	"github.com/hootuu/hyper/payment"
	"github.com/hootuu/hyper/shipping"
	"time"
)

type Dealer struct {
	code    hiorder.Code
	timeout time.Duration

	f     *hiorder.Factory[Matter]
	topup *TopUp
}

func newDealer(
	code hiorder.Code,
	timeout time.Duration,
) *Dealer {
	return &Dealer{
		code:    code,
		timeout: timeout,
	}
}

func (d *Dealer) doInit(f *hiorder.Factory[Matter], topup *TopUp) {
	d.topup = topup
	d.f = f
}

func (d *Dealer) Code() hiorder.Code {
	return d.code
}

func (d *Dealer) Build(ord hiorder.Order[Matter]) (hiorder.Deal[Matter], error) {
	return newDeal(&ord, d, d.topup), nil
}

// OnPaymentAltered todo use ctx for mq
func (d *Dealer) OnPaymentAltered(ctx context.Context, alter *payment.AlterPayload) error {

	if alter == nil {
		hlog.Fix("hitopup.dealer.OnPaymentAltered: alter is nil")
		return nil
	}
	//todo add statuc check
	//ctx := context.Background()
	//_, err := d.f.Load(ctx, cast.ToUint64(alter.BizID)) //todo
	//if err != nil {
	//	hlog.Err("hitopup.dealer.OnPaymentAltered: load fail", zap.Error(err))
	//	return err
	//}
	//fmt.Println("hitopup.dealer.OnPaymentAltered: load success") // todo
	////err = eng.
	////if err != nil {
	////	hlog.Err("hitopup.dealer.OnPaymentAltered: complete fail", zap.Error(err))
	////	return err
	////}todo
	return nil
}

func (d *Dealer) OnShippingAltered(ctx context.Context, alter *shipping.AlterPayload) error {
	return nil
}
