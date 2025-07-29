package payment

import (
	"context"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/hfsm"
	"github.com/hootuu/hyle/hlog"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

func advJobToPrepared(
	ctx context.Context,
	current hfsm.State,
	_ hfsm.Event,
	data ...dict.Dict,
) (hfsm.State, error) {
	jobID := cast.ToString(data[0].Get("job_id"))
	err := doSetJobStatus(ctx, jobID, JobPrepared, current, nil)
	if err != nil {
		return 0, err
	}
	return JobPrepared, nil
}

func doAdvJobToPrepared(ctx context.Context, pid ID, seq int) (err error) {
	if hlog.IsElapseComponent() {
		defer hlog.ElapseWithCtx(ctx,
			"hyper.payment.doAdvJobToPrepared",
			hlog.F(zap.Uint64("pid", pid), zap.Int("seq", seq)),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err)}
				}
				return nil
			},
		)()
	}

	jobID := BuildJobID(pid, seq)
	err = doJobAdvance(ctx, jobID, JobPrepareEvent, nil)
	if err != nil {
		return err
	}
	return nil
}
