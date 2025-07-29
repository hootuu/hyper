package payment

import (
	"context"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/hypes/collar"
	"time"
)

type ID = uint64
type Payment struct {
	ID       ID            `json:"id"`
	Payer    collar.Link   `json:"payer"`
	Payee    collar.Link   `json:"payee"`
	BizCode  string        `json:"biz_code"`
	BizID    string        `json:"biz_id"`
	Amount   uint64        `json:"amount"`
	Status   Status        `json:"status"`
	JobCount int           `json:"job_count"`
	Timeout  time.Duration `json:"timeout"`
}

type Channel = string
type JobID = string
type Job struct {
	JobID      JobID         `json:"job_id"`
	Channel    Channel       `json:"channel"`
	PaymentID  ID            `json:"payment_id"`
	PaymentSeq int           `json:"payment_seq"`
	Status     JobStatus     `json:"status"`
	Ctx        dict.Dict     `json:"ctx"`
	Timeout    time.Duration `json:"timeout"`
}

type JobDefine interface {
	Validate() error
	GetChannel() Channel
	GetAmount() uint64
	GetCtx() dict.Dict
}

type JobExecutor interface {
	GetChannel() Channel
	Prepare(ctx context.Context, pay *Payment, job *Job) (synced bool, err error)
	Advance(ctx context.Context, pay *Payment, job *Job) (synced bool, err error)
	Cancel(ctx context.Context, pay *Payment, job *Job) (synced bool, err error)
	Timeout(ctx context.Context, pay *Payment, job *Job) (synced bool, err error)
	OnPrepared(ctx context.Context, job *Job) error
	OnCompleted(ctx context.Context, job *Job) error
	OnTimeout(ctx context.Context, job *Job) error
	OnCanceled(ctx context.Context, job *Job) error
}
