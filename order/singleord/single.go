package singleord

import (
	"context"
	"github.com/hootuu/hyle/data/ctrl"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/tag"
	"github.com/hootuu/hyle/hcoin"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyper/hiorder"
	"github.com/hootuu/hyper/hiprod/prod"
	"github.com/spf13/cast"
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
	SkuID        prod.SkuID   `json:"sku_id"`
	Payer        collar.Link  `json:"payer"`
	PayerAccount collar.Link  `json:"payer_account"`
	Quantity     uint32       `json:"quantity"`
	Amount       hcoin.Amount `json:"amount"`
	Ctrl         ctrl.Ctrl    `json:"ctrl"`
	Tag          tag.Tag      `json:"tag"`
	Meta         dict.Dict    `json:"meta"`
}

func (s *Single) Create(ctx context.Context, paras *CreateParas) (*hiorder.Order[Matter], error) {
	engine, err := s.factory.New(ctx, &hiorder.CreateParas[Matter]{
		Title:        "[PROD]" + cast.ToString(paras.SkuID),
		Payer:        paras.Payer,
		PayerAccount: paras.PayerAccount,
		Payee:        s.payee,
		PayeeAccount: s.payee,
		Amount:       paras.Amount,
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

func (s *Single) doInit() error {
	d := newDealer(s.code, hcoin.Currency("CNY"), 15*time.Minute)
	s.factory = hiorder.NewFactory[Matter](d)
	d.doInit(s.factory, s)
	return nil
}
