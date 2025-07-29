package payment

import (
	"context"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/hfsm"
	"github.com/spf13/cast"
)

func advJobToCanceled(
	ctx context.Context,
	current hfsm.State,
	_ hfsm.Event,
	data ...dict.Dict,
) (hfsm.State, error) {
	jobID := cast.ToString(data[0].Get("job_id"))
	err := doSetJobStatus(ctx, jobID, JobCanceled, current, nil)
	if err != nil {
		return 0, err
	}
	return JobCanceled, nil
}
