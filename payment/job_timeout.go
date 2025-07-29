package payment

import (
	"context"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/hfsm"
	"github.com/hootuu/hyper/hyperplt"
	"github.com/spf13/cast"
)

func advJobToTimeout(
	ctx context.Context,
	current hfsm.State,
	_ hfsm.Event,
	data ...dict.Dict,
) (hfsm.State, error) {
	jobID := cast.ToString(data[0].Get("job_id"))
	err := doSetJobStatus(ctx, jobID, JobTimeout, current, nil)
	if err != nil {
		return 0, err
	}
	return JobTimeout, nil
}

func callbackJobCheckTimeoutAndAdv(ctx context.Context, jobID JobID) error {
	jobM, err := hdb.Get[JobM](hyperplt.Tx(ctx), "id = ?", jobID)
	if err != nil {
		return err
	}
	if jobM == nil {
		return nil
	}
	switch jobM.Status {
	case JobTimeout:
		return nil
	case JobPrepared:
		return doJobAdvance(ctx, jobID, JobTimeoutEvent, nil)
	default:
		return nil
	}
}
