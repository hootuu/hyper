package thirdjob

import (
	"context"
	"github.com/hootuu/hyper/hpay/payment"
)

type Executor struct {
}

func NewExecutor() *Executor {
	return &Executor{}
}

func (e *Executor) GetChannel() payment.Channel {
	return ThirdChannel
}

func (e *Executor) Prepare(_ context.Context, _ *payment.Payment, _ *payment.Job) (synced bool, err error) {
	return false, nil
}

func (e *Executor) Advance(_ context.Context, _ *payment.Payment, _ *payment.Job) (synced bool, err error) {
	return false, nil
}

func (e *Executor) Cancel(ctx context.Context, job *payment.Job) error {
	//TODO implement me
	panic("implement me")
}

func (e *Executor) OnPrepared(ctx context.Context, job *payment.Job) error {
	return payment.OnJobPrepared(ctx, job)
}

func (e *Executor) OnCompleted(ctx context.Context, job *payment.Job) error {
	return payment.OnJobCompleted(ctx, job)
}

func (e *Executor) OnTimeout(ctx context.Context, job *payment.Job) error {
	return payment.OnJobTimeout(ctx, job)
}

func (e *Executor) OnCanceled(ctx context.Context, job *payment.Job) error {
	return payment.OnJobCanceled(ctx, job)
}
