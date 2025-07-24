package singleord

import (
	"context"
	"errors"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/ctrl"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/tag"
	"github.com/hootuu/hyle/hcoin"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyle/hypes/ex"
	"github.com/hootuu/hyper/hiorder"
	"github.com/hootuu/hyper/hiprod/prod"
	"github.com/hootuu/hyper/hpay"
	"github.com/hootuu/hyper/hpay/payment"
	"github.com/hootuu/hyper/hpay/thirdjob"
	"github.com/hootuu/hyper/hshipping"
	"github.com/hootuu/hyper/hshipping/shipping"
	"github.com/hootuu/hyper/hyperplt"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type Single struct {
	payee   collar.Link
	code    hiorder.Code
	factory *hiorder.Factory[Matter]
}

func Build(code hiorder.Code, payee collar.Link) (*Single, error) {
	s := &Single{
		code:  code,
		payee: payee,
	}
	err := s.doInit()
	if err != nil {
		return nil, err
	}
	return s, nil
}

type CreateParas struct {
	SkuID    prod.SkuID   `json:"sku_id"`
	Payer    collar.Link  `json:"payer"`
	Quantity uint32       `json:"quantity"`
	Amount   hcoin.Amount `json:"amount"`
	Ctrl     ctrl.Ctrl    `json:"ctrl"`
	Tag      tag.Tag      `json:"tag"`
	Meta     dict.Dict    `json:"meta"`
	UniLink  collar.Link  `json:"uniLink"`
}

func (s *Single) Create(ctx context.Context, paras *CreateParas) (*hiorder.Order[Matter], error) {
	engine, err := s.factory.New(ctx, &hiorder.CreateParas[Matter]{
		Title:   "[PROD]" + cast.ToString(paras.SkuID),
		Payer:   paras.Payer,
		Payee:   s.payee,
		Amount:  paras.Amount,
		Payment: nil,
		Matter: Matter{
			SkuID:    paras.SkuID,
			VwhID:    0,
			PwhID:    0,
			Price:    1000,
			Quantity: paras.Quantity,
			Amount:   paras.Amount,
		},
		Ctrl:    paras.Ctrl,
		Tag:     paras.Tag,
		Meta:    paras.Meta,
		UniLink: paras.UniLink,
	})
	if err != nil {
		return nil, err
	}
	err = engine.Submit(ctx)
	if err != nil {
		return nil, err
	}
	return engine.GetOrder(), nil
}

func (s *Single) PaymentPrepared(
	ctx context.Context,
	ordID hiorder.ID,
	chanCode string,
	exM *ex.Ex,
) (err error) {
	if hlog.IsElapseComponent() {
		defer hlog.ElapseWithCtx(ctx, "hyper.single.PaymentPrepared",
			hlog.F(zap.Uint64("ordID", ordID)),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err)}
				}
				return nil
			})()
	}
	eng, err := s.factory.Load(ctx, ordID)
	if err != nil {
		hlog.Err("s.factory.Load", hlog.TraceInfo(ctx), zap.Error(err))
		return errors.New("load order fail: " + err.Error())
	}
	ord := eng.GetOrder()
	paymentID, err := s.factory.SetPayment(ctx, ord.ID, []payment.JobDefine{
		&thirdjob.Job{
			ThirdCode: chanCode,
			Amount:    ord.Amount,
			Ex:        exM,
		}},
		exM,
	)
	if err != nil {
		hlog.Err("s.factory.SetPayment", hlog.TraceInfo(ctx), zap.Error(err))
		return errors.New("SetPayment fail: " + err.Error())
	}

	err = hpay.JobPrepared(ctx, paymentID, 1)
	if err != nil {
		return errors.New("job prepared: " + err.Error())
	}
	err = hpay.Advance(ctx, paymentID)
	if err != nil {
		return errors.New("advance payment: " + err.Error())
	}
	return nil
}

func (s *Single) GetIdByUniLink(uniLink collar.Link) (hiorder.ID, error) {
	return s.factory.GetIDByUniLink(uniLink)
}

func (s *Single) PaymentCompleted(
	ctx context.Context,
	ordID hiorder.ID,
	payNumber string,
) (err error) {
	if hlog.IsElapseComponent() {
		defer hlog.ElapseWithCtx(ctx, "hyper.topup.TopUpPaymentCompleted",
			hlog.F(zap.Uint64("ordID", ordID)),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err)}
				}
				return nil
			})()
	}
	eng, err := s.factory.Load(ctx, ordID)
	if err != nil {
		return errors.New("load order fail: " + err.Error())
	}
	paymentID := eng.GetOrder().PaymentID
	err = hpay.JobCompleted(ctx, paymentID, 1, payNumber)
	if err != nil {
		return errors.New("job completed: " + err.Error())
	}
	return nil
}

func (s *Single) ShippingCreate(
	ctx context.Context,
	ordID hiorder.ID,
	addr *shipping.Address,
	ex *ex.Ex,
) (err error) {
	if hlog.IsElapseComponent() {
		defer hlog.ElapseWithCtx(ctx, "hyper.single.ShippingPrepared",
			hlog.F(zap.Uint64("ordID", ordID)),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err)}
				}
				return nil
			})()
	}
	eng, err := s.factory.Load(ctx, ordID)
	if err != nil {
		return errors.New("load order fail: " + err.Error())
	}
	uniLink := s.BuildShippingCollar(ordID)
	err = hdb.Tx(hyperplt.Tx(ctx), func(tx *gorm.DB) error {
		innerCtx := hdb.TxCtx(tx)
		shippingID, err := hshipping.ShippingCreate(innerCtx, hshipping.CreateParas{
			UniLink: uniLink.Link(),
			Address: addr,
			Ex:      ex,
		})
		if err != nil {
			return err
		}
		err = eng.SetShipping(innerCtx, ordID, shippingID)
		if err != nil {
			return err
		}
		return nil
	})

	return nil
}

func (s *Single) ShippingPrepared(
	ctx context.Context,
	ordID hiorder.ID,
	courierCode shipping.CourierCode,
	trackingNo string,
	meta ex.Meta,
) (err error) {
	if hlog.IsElapseComponent() {
		defer hlog.ElapseWithCtx(ctx, "hyper.single.ShippingPrepared",
			hlog.F(zap.Uint64("ordID", ordID)),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err)}
				}
				return nil
			})()
	}
	uniLink := s.BuildShippingCollar(ordID)
	err = hshipping.ShippingPrepared(
		ctx,
		uniLink.Link(),
		courierCode,
		trackingNo,
		meta,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Single) ShippingCompleted(ctx context.Context, ordID hiorder.ID, meta ex.Meta) (err error) {
	if hlog.IsElapseComponent() {
		defer hlog.ElapseWithCtx(ctx, "hyper.single.ShippingCompleted",
			hlog.F(zap.Uint64("ordID", ordID)),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err)}
				}
				return nil
			})()
	}
	uniLink := s.BuildShippingCollar(ordID)
	err = hshipping.ShippingCompleted(ctx, uniLink.Link(), meta)
	if err != nil {
		return err
	}
	return nil
}

func (s *Single) BuildShippingCollar(ordID hiorder.ID) collar.Collar {
	return collar.Build(s.code, cast.ToString(ordID))
}

func (s *Single) doInit() error {
	d := newDealer(s.code, 15*time.Minute)
	s.factory = hiorder.NewFactory[Matter](d)
	d.doInit(s.factory, s)
	return nil
}
