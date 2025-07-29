package payment

import "context"

type ThirdExecutor struct {
}

func NewThirdExecutor() *ThirdExecutor {
	return &ThirdExecutor{}
}

func (e *ThirdExecutor) GetChannel() Channel {
	return ThirdChannel
}

func (e *ThirdExecutor) Prepare(ctx context.Context, pay *Payment, job *Job) (synced bool, err error) {

	return false, nil
}

func (e *ThirdExecutor) Advance(_ context.Context, _ *Payment, _ *Job) (synced bool, err error) {
	return false, nil
}

func (e *ThirdExecutor) Cancel(ctx context.Context, pay *Payment, job *Job) (synced bool, err error) {
	return false, nil
}

func (e *ThirdExecutor) Timeout(ctx context.Context, pay *Payment, job *Job) (synced bool, err error) {
	return false, nil
}

func (e *ThirdExecutor) OnPrepared(ctx context.Context, job *Job) error {
	return nil
}

func (e *ThirdExecutor) OnCompleted(ctx context.Context, job *Job) error {
	return nil
}

func (e *ThirdExecutor) OnTimeout(ctx context.Context, job *Job) error {
	return nil
}

func (e *ThirdExecutor) OnCanceled(ctx context.Context, job *Job) error {
	return nil
}
