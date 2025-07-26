package singleord

import (
	"context"
	"errors"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/hiorder"
	"github.com/hootuu/hyper/hpay/payment"
	"github.com/hootuu/hyper/hshipping/shipping"
	"go.uber.org/zap"
	"time"
)

type Dealer struct {
	code    hiorder.Code
	timeout time.Duration

	f      *hiorder.Factory[Matter]
	single *Single
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

func (d *Dealer) Code() hiorder.Code {
	return d.code
}

func (d *Dealer) Build(ord hiorder.Order[Matter]) (hiorder.Deal[Matter], error) {
	return newDeal(d, &ord), nil
}

func (d *Dealer) OnPaymentAltered(alter *hiorder.PaymentAltered[Matter]) (err error) {
	ctx := context.Background() //todo add ctx log
	if hlog.IsElapseFunction() {
		defer hlog.ElapseWithCtx(ctx, "hyper.singleord.OnPaymentAltered",
			hlog.F(zap.Uint64("payment", alter.PaymentID), zap.Uint64("ord", alter.Order.ID)),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err)}
				}
				return nil
			},
		)()
	}
	switch alter.DstStatus {
	case payment.Completed:
		//todo fix logs
		eng, err := d.f.Load(ctx, alter.Order.ID)
		if err != nil {
			hlog.Err("load engine failed", hlog.TraceInfo(ctx), zap.Error(err))
			return errors.New("load engine failed: " + err.Error())
		}
		err = eng.Consense(ctx)
		if err != nil {
			hlog.Err("complete engine failed", hlog.TraceInfo(ctx), zap.Error(err))
			return errors.New("complete engine failed: " + err.Error())
		}
	default:
		return nil
	}
	return nil
}

func (d *Dealer) OnShippingAltered(alter *hiorder.ShippingAltered[Matter]) (err error) {
	ctx := context.Background() //todo add ctx log
	if hlog.IsElapseFunction() {
		defer hlog.ElapseWithCtx(ctx, "hyper.singleord.OnShippingAltered",
			hlog.F(zap.Uint64("shippingID", alter.ShippingID), zap.Uint64("ord", alter.Order.ID)),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err)}
				}
				return nil
			},
		)()
	}
	switch alter.DstStatus {
	case shipping.StatusPickedUp:
		//todo fix logs
		eng, err := d.f.Load(ctx, alter.Order.ID)
		if err != nil {
			hlog.Err("load engine failed", hlog.TraceInfo(ctx), zap.Error(err))
			return errors.New("load engine failed: " + err.Error())
		}
		err = eng.Execute(ctx)
		if err != nil {
			hlog.Err("complete engine failed", hlog.TraceInfo(ctx), zap.Error(err))
			return errors.New("complete engine failed: " + err.Error())
		}
	case shipping.StatusDelivered:
		eng, err := d.f.Load(ctx, alter.Order.ID)
		if err != nil {
			hlog.Err("load engine failed", hlog.TraceInfo(ctx), zap.Error(err))
			return errors.New("load engine failed: " + err.Error())
		}
		err = eng.Complete(ctx)
		if err != nil {
			hlog.Err("complete engine failed", hlog.TraceInfo(ctx), zap.Error(err))
			return errors.New("complete engine failed: " + err.Error())
		}
	default:
		return nil
	}
	return nil
}

func (d *Dealer) doInit(f *hiorder.Factory[Matter], single *Single) {
	d.single = single
	d.f = f
}
