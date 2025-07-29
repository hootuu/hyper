package payment

import (
	"errors"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/hypes/ex"
	"github.com/nineora/harmonic/chain"
	"github.com/spf13/cast"
)

type NineJob struct {
	Mint   chain.Address `json:"mint"`
	Payer  chain.Address `json:"payer"`
	Payee  chain.Address `json:"payee"`
	Amount uint64        `json:"amount"`
	Ex     *ex.Ex        `json:"ex"`
}

func (j *NineJob) Validate() error {
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
	return nil
}

func (j *NineJob) GetChannel() Channel {
	return NineChannel
}

func (j *NineJob) GetCtx() dict.Dict {
	return dict.New(map[string]any{
		"mint":   j.Mint,
		"payer":  j.Payer,
		"payee":  j.Payee,
		"amount": j.Amount,
	})
}

func (j *NineJob) GetAmount() uint64 {
	return j.Amount
}

func nineJobFromCtx(ctx dict.Dict) (*NineJob, error) {
	return &NineJob{
		Mint:   ctx.Get("mint").String(),
		Payer:  ctx.Get("payer").String(),
		Payee:  ctx.Get("payee").String(),
		Amount: cast.ToUint64(ctx.Get("amount").String()),
		//Ex:     dict.NewDict(ctx.Get("ex").Data()), todo add fix ex
	}, nil
}
