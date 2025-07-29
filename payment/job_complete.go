package payment

import (
	"context"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/hfsm"
	"github.com/spf13/cast"
)

func advJobToCompleted(
	ctx context.Context,
	current hfsm.State,
	_ hfsm.Event,
	data ...dict.Dict,
) (hfsm.State, error) {
	jobID := cast.ToString(data[0].Get("job_id"))
	payNo := data[0].Get("pay_no").String()
	err := doSetJobStatus(ctx, jobID, JobCompleted, current, dict.NewDict().Set("pay_no", payNo))
	if err != nil {
		return 0, err
	}
	return JobCompleted, nil
}

func doAdvJobToCompleted(ctx context.Context, pid ID, seq int, payNo string) (err error) {
	jobID := BuildJobID(pid, seq)
	return doJobAdvance(ctx, jobID, JobCompleteEvent, dict.NewDict().Set("pay_no", payNo))
}
