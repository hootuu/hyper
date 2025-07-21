package hitopup

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/unicom/hmq/hmq"
	"github.com/hootuu/hyle/data/ctrl"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/tag"
	"github.com/hootuu/hyle/hcfg"
	"github.com/hootuu/hyle/hcoin"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyper/hiorder"
	"github.com/nineora/harmonic/chain"
	"strings"
	"time"
)

type Convert func(src uint64) uint64

type TopUp struct {
	code       hiorder.Code
	mint       chain.Address
	factory    *hiorder.Factory[Matter]
	convert    Convert
	mqConsumer *hmq.Consumer
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
	Title        string       `json:"title"`
	Payer        collar.Link  `json:"payer"`
	PayerAccount collar.Link  `json:"payer_account"`
	Amount       hcoin.Amount `json:"amount"`
	Ctrl         ctrl.Ctrl    `json:"ctrl"`
	Tag          tag.Tag      `json:"tag"`
	Meta         dict.Dict    `json:"meta"`
}

func (p TopUpParas) validate() error {
	return nil //todo
}

func (t *TopUp) TopUp(ctx context.Context, paras TopUpParas) (*hiorder.Order[Matter], error) {
	if err := paras.validate(); err != nil {
		return nil, err
	}
	engine, err := t.factory.New(ctx, &hiorder.CreateParas[Matter]{
		Title:        paras.Title,
		Payer:        paras.Payer,
		PayerAccount: paras.PayerAccount,
		Payee:        collar.Build("NINEORA", t.mint).Link(),
		PayeeAccount: collar.Build("NINEORA", t.mint).Link(), //todo
		Amount:       paras.Amount,
		Matter:       Matter{},
		Ctrl:         paras.Ctrl,
		Tag:          paras.Tag,
		Meta:         paras.Meta,
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

func (t *TopUp) doInit() error {
	currencyStr := hcfg.GetString(t.cfg("currency"), string(hcoin.CNY))
	timeout := hcfg.GetDuration(t.cfg("timeout"), 15*time.Minute)
	d := newDealer(t.code, hcoin.Currency(currencyStr), timeout)
	t.factory = hiorder.NewFactory[Matter](d)
	d.doInit(t.factory, t)
	return nil
}

func (t *TopUp) cfg(k string) string {
	return fmt.Sprintf("hitopup.%s.%s", strings.ToLower(string(t.code)), k)
}
