package hitopup

import (
	"context"
	"errors"
	"fmt"
	"github.com/hootuu/hyle/data/ctrl"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/tag"
	"github.com/hootuu/hyle/hcfg"
	"github.com/hootuu/hyle/hcoin"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyle/hypes/ex"
	"github.com/hootuu/hyper/hiorder"
	"github.com/hootuu/hyper/hpay"
	"github.com/hootuu/hyper/hpay/payment"
	"github.com/hootuu/hyper/hpay/thirdjob"
	"github.com/nineora/harmonic/chain"
	"go.uber.org/zap"
	"strings"
	"time"
)

type Convert func(src uint64) uint64

type TopUp struct {
	code    hiorder.Code
	mint    chain.Address
	factory *hiorder.Factory[Matter]
	convert Convert
}

func NewTopUp(code hiorder.Code, mint chain.Address, convert Convert) (*TopUp, error) {
	t := &TopUp{
		code: code,
		mint: mint,
	}
	if convert != nil {
		t.convert = convert
	} else {
		t.convert = func(src uint64) uint64 {
			return src
		}
	}
	err := t.doInit()
	if err != nil {
		return nil, err
	}
	return t, nil
}

type TopUpParas struct {
	Title          string        `json:"title"`
	Payer          collar.Link   `json:"payer"`
	Amount         hcoin.Amount  `json:"amount"`
	PayChannelCode string        `json:"pay_channel_code"`
	InAccountAddr  chain.Address `json:"in_account_addr"`
	Ctrl           ctrl.Ctrl     `json:"ctrl"`
	Tag            tag.Tag       `json:"tag"`
	Meta           dict.Dict     `json:"meta"`
}

func (p TopUpParas) validate() error {
	if p.Title == "" {
		return errors.New("title is required")
	}
	if p.Payer == "" {
		return errors.New("payer is required")
	}
	if p.Amount == 0 {
		return errors.New("amount is required")
	}
	if p.PayChannelCode == "" {
		return errors.New("pay_channel_code is required")
	}
	if p.InAccountAddr == "" {
		return errors.New("in_account_addr is required")
	}
	return nil
}

func (t *TopUp) TopUpCreate(ctx context.Context, paras TopUpParas) (*hiorder.Order[Matter], error) {
	if err := paras.validate(); err != nil {
		return nil, err
	}
	engine, err := t.factory.New(ctx, &hiorder.CreateParas[Matter]{
		Title:  paras.Title,
		Payer:  paras.Payer,
		Payee:  collar.Build("NINEORA", t.mint).Link(),
		Amount: paras.Amount,
		Payment: []payment.JobDefine{&thirdjob.Job{
			ThirdCode: paras.PayChannelCode,
			Amount:    paras.Amount,
			Ex: &ex.Ex{
				Ctrl: paras.Ctrl,
				Tag:  paras.Tag,
				Meta: paras.Meta,
			},
		}},
		Matter: Matter{InAccount: paras.InAccountAddr},
		Ctrl:   paras.Ctrl,
		Tag:    paras.Tag,
		Meta:   paras.Meta,
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

func (t *TopUp) TopUpPaymentPrepared(ctx context.Context, ordID hiorder.ID) (err error) {
	if hlog.IsElapseComponent() {
		defer hlog.ElapseWithCtx(ctx, "hyper.topup.TopUpPaymentPrepared",
			hlog.F(zap.Uint64("ordID", ordID)),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err)}
				}
				return nil
			})()
	}
	eng, err := t.factory.Load(ctx, ordID)
	if err != nil {
		return errors.New("load order fail: " + err.Error())
	}
	paymentID := eng.GetOrder().PaymentID
	err = hpay.JobPrepared(ctx, paymentID, 1, t.code)
	if err != nil {
		return errors.New("job prepared: " + err.Error())
	}
	err = hpay.Advance(ctx, paymentID)
	if err != nil {
		return errors.New("advance payment: " + err.Error())
	}
	return nil
}

func (t *TopUp) TopUpPaymentCompleted(ctx context.Context, ordID hiorder.ID, payNumber string) (err error) {
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
	eng, err := t.factory.Load(ctx, ordID)
	if err != nil {
		return errors.New("load order fail: " + err.Error())
	}
	paymentID := eng.GetOrder().PaymentID
	err = hpay.JobCompleted(ctx, paymentID, 1, eng.GetOrder().Code, payNumber)
	if err != nil {
		return errors.New("job completed: " + err.Error())
	}
	return nil
}

func (t *TopUp) doInit() error {
	timeout := hcfg.GetDuration(t.cfg("timeout"), 15*time.Minute)
	d := newDealer(t.code, timeout)
	t.factory = hiorder.NewFactory[Matter](d)
	d.doInit(t.factory, t)
	return nil
}

func (t *TopUp) cfg(k string) string {
	return fmt.Sprintf("hitopup.%s.%s", strings.ToLower(string(t.code)), k)
}
