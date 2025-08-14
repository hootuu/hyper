package payment

import (
	"errors"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/hypes/ex"
)

const (
	ThirdChannel = "THIRD"
)

type ThirdJob struct {
	ThirdCode string `json:"third_code"`
	Amount    uint64 `json:"amount"`
	Ex        *ex.Ex `json:"ex"`
}

func (j *ThirdJob) Validate() error {
	if j.ThirdCode == "" {
		return errors.New("require third_code")
	}
	//if j.Amount == 0 {
	//	return errors.New("require amount")
	//}
	return nil
}

func (j *ThirdJob) GetChannel() Channel {
	return ThirdChannel
}

func (j *ThirdJob) GetAmount() uint64 {
	return j.Amount
}

func (j *ThirdJob) GetCtx() dict.Dict {
	return dict.New(map[string]any{
		"third_code": j.ThirdCode,
		"amount":     j.Amount,
		//todo add ex
	})
}
