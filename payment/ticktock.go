package payment

import (
	"context"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/ticktock"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hretry"
	"github.com/hootuu/hyper/hyperplt"
	"go.uber.org/zap"
)

const (
	ttPaymentJobTimeout = "HYPER_PAYMENT_JOB_TIMEOUT"
)

func ttListenJobTimeout(ctx context.Context, job *Job) {
	if job.Timeout == 0 {
		return
	}
	hretry.Universal(func() error {
		err := hyperplt.Postman().Send(ctx, &ticktock.DelayJob{
			Type:      ttPaymentJobTimeout,
			ID:        job.JobID,
			Payload:   []byte(job.JobID),
			UniqueTTL: 0,
			Delay:     job.Timeout,
		})
		if err != nil {
			hlog.TraceFix("send ticktock failed", ctx, err, zap.Any("job", job))
			return err
		}
		return err
	})
}

//todo use ticktock init

func init() {
	helix.Ready(func() {
		hyperplt.Ticktock().RegisterJobHandlerFunc(
			ttPaymentJobTimeout,
			func(ctx context.Context, job *ticktock.Job) error {
				return callbackJobCheckTimeoutAndAdv(ctx, string(job.Payload))
			},
		)
	})
}
