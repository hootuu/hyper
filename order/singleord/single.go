package singleord

import (
	"context"
	"errors"
	"github.com/hootuu/hyle/data/ctrl"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/tag"
	"github.com/hootuu/hyle/hcoin"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyper/hiorder"
	"github.com/hootuu/hyper/hiprod/prod"
	"github.com/hootuu/hyper/hpay"
	"github.com/hootuu/hyper/hpay/payment"
	"github.com/spf13/cast"
	"go.uber.org/zap"
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
	SkuID    prod.SkuID          `json:"sku_id"`
	Payer    collar.Link         `json:"payer"`
	Quantity uint32              `json:"quantity"`
	Amount   hcoin.Amount        `json:"amount"`
	Payment  []payment.JobDefine `json:"payment"`
	Ctrl     ctrl.Ctrl           `json:"ctrl"`
	Tag      tag.Tag             `json:"tag"`
	Meta     dict.Dict           `json:"meta"`
}

func (s *Single) Create(ctx context.Context, paras *CreateParas) (*hiorder.Order[Matter], error) {
	engine, err := s.factory.New(ctx, &hiorder.CreateParas[Matter]{
		Title:   "[PROD]" + cast.ToString(paras.SkuID),
		Payer:   paras.Payer,
		Payee:   s.payee,
		Amount:  paras.Amount,
		Payment: paras.Payment,
		Matter: Matter{
			SkuID:    paras.SkuID,
			VwhID:    0,
			PwhID:    0,
			Price:    1000,
			Quantity: paras.Quantity,
			Amount:   paras.Amount,
		},
		Ctrl: paras.Ctrl,
		Tag:  paras.Tag,
		Meta: paras.Meta,
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

func (s *Single) PaymentPrepared(ctx context.Context, ordID hiorder.ID) (err error) {
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
		return errors.New("load order fail: " + err.Error())
	}
	paymentID := eng.GetOrder().PaymentID
	err = hpay.JobPrepared(ctx, paymentID, 1, s.code)
	if err != nil {
		return errors.New("job prepared: " + err.Error())
	}
	err = hpay.Advance(ctx, paymentID)
	if err != nil {
		return errors.New("advance payment: " + err.Error())
	}
	return nil
}

func (s *Single) TopUpPaymentCompleted(ctx context.Context, ordID hiorder.ID) (err error) {
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
	err = hpay.JobCompleted(ctx, paymentID, 1, s.code)
	if err != nil {
		return errors.New("job completed: " + err.Error())
	}
	return nil
}

func (s *Single) doInit() error {
	d := newDealer(s.code, 15*time.Minute)
	s.factory = hiorder.NewFactory[Matter](d)
	d.doInit(s.factory, s)
	return nil
}
