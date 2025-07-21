package ninejob

import (
	"errors"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/idx"
	"github.com/hootuu/hyle/hypes/ex"
	"github.com/hootuu/hyper/hpay/payment"
	"github.com/nineora/harmonic/chain"
	"github.com/spf13/cast"
)

type Job struct {
	Mint      chain.Address `json:"mint"`
	Payer     chain.Address `json:"payer"`
	Payee     chain.Address `json:"payee"`
	Amount    uint64        `json:"amount"`
	Ex        *ex.Ex        `json:"ex"`
	CheckCode string        `json:"check_code"`
}

func (j *Job) Validate() error {
	if j.Mint == "" {
		return errors.New("mint address is required")
	}
	if j.Payer == "" {
		return errors.New("payer address is required")
	}
	if j.Payee == "" {
		return errors.New("payee address is required")
	}
	if j.Amount == 0 {
		return errors.New("amount is required")
	}
	if j.CheckCode == "" {
		j.CheckCode = idx.New()
	}
	return nil
}

func (j *Job) GetChannel() payment.Channel {
	return NineChannel
}

func (j *Job) GetCtx() dict.Dict {
	return dict.New(map[string]any{
		"mint":   j.Mint,
		"payer":  j.Payer,
		"payee":  j.Payee,
		"amount": j.Amount,
	})
}

func (j *Job) GetAmount() uint64 {
	return j.Amount
}

func (j *Job) GetCheckCode() string {
	return j.CheckCode
}

func JobFromCtx(ctx dict.Dict) (*Job, error) {
	return &Job{
		Mint:   ctx.Get("mint").String(),
		Payer:  ctx.Get("payer").String(),
		Payee:  ctx.Get("payee").String(),
		Amount: cast.ToUint64(ctx.Get("amount").String()),
		//Ex:     dict.NewDict(ctx.Get("ex").Data()), todo add fix ex
		CheckCode: ctx.Get("check_code").String(),
	}, nil
}
