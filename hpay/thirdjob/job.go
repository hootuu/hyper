package thirdjob

import (
	"errors"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/hypes/ex"
	"github.com/hootuu/hyper/hpay/payment"
)

const (
	ThirdChannel = "THIRD"
)

type Job struct {
	ThirdCode string `json:"third_code"`
	Amount    uint64 `json:"amount"`
	Ex        *ex.Ex `json:"ex"`
}

func (j *Job) Validate() error {
	if j.ThirdCode == "" {
		return errors.New("require third_code")
	}
	if j.Amount == 0 {
		return errors.New("require amount")
	}
	return nil
}

func (j *Job) GetChannel() payment.Channel {
	return ThirdChannel
}

func (j *Job) GetAmount() uint64 {
	return j.Amount
}

func (j *Job) GetCtx() dict.Dict {
	return dict.New(map[string]any{
		"third_code": j.ThirdCode,
		"amount":     j.Amount,
		//todo add ex
	})
}
