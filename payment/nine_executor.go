package payment

import (
	"context"
	"errors"
	"fmt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyper/hyperplt"
	"github.com/nineora/harmonic/harmonic"
	"github.com/nineora/harmonic/nineapi"
	"github.com/nineora/harmonic/nineora"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	NineChannel = "NINEORA"
)

type NineExecutor struct{}

func NewNineExecutor() *NineExecutor {
	return &NineExecutor{}
}

func (e *NineExecutor) GetChannel() Channel {
	return NineChannel
}

func (e *NineExecutor) Prepare(ctx context.Context, pay *Payment, job *Job) (synced bool, err error) {
	if pay == nil {
		return false, fmt.Errorf("assert: pay is not nil")
	}
	if job == nil {
		return false, fmt.Errorf("assert: job is not nil")
	}
	if pay.ID != job.PaymentID {
		return false, fmt.Errorf("assert: payment from job %d is not %d", job.PaymentID, pay.ID)
	}
	if hlog.IsElapseComponent() {
		defer hlog.ElapseWithCtx(ctx, "hyper.payment.nine.Prepare",
			hlog.F(zap.String("jobID", job.JobID)),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err)}
				}
				return nil
			},
		)()
	}

	localJob, err := nineJobFromCtx(job.Ctx)
	if err != nil {
		return false, err
	}

	_, err = harmonic.Nineora().TokenLock(ctx, &nineapi.TxLockParas{
		Account: localJob.Payer,
		Amount:  localJob.Amount,
		Biz:     pay.BizCode,
		Link:    collar.Build(pay.BizCode, pay.BizID).Link(),
		Ex:      localJob.Ex,
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (e *NineExecutor) Advance(ctx context.Context, pay *Payment, job *Job) (synced bool, err error) {
	if pay == nil {
		return false, fmt.Errorf("assert: pay is not nil")
	}
	if job == nil {
		return false, fmt.Errorf("assert: job is not nil")
	}
	if pay.ID != job.PaymentID {
		return false, fmt.Errorf("assert: payment from job %d is not %d", job.PaymentID, pay.ID)
	}
	var unlockSign nineora.Signature = ""
	var transSign nineora.Signature = ""
	if hlog.IsElapseComponent() {
		defer hlog.ElapseWithCtx(ctx, "hyper.payment.nine.Advance",
			hlog.F(zap.String("jobID", job.JobID)),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err)}
				}
				return []zap.Field{
					zap.String("unlockSign", string(unlockSign)),
					zap.String("transSign", string(transSign)),
				}
			},
		)()
	}

	localJob, err := nineJobFromCtx(job.Ctx)
	if err != nil {
		return false, err
	}

	bizID := collar.Build(pay.BizCode, pay.BizID).Link()
	err = hdb.Tx(hyperplt.Tx(ctx), func(tx *gorm.DB) error {
		innerTx := hdb.TxCtx(tx, ctx)
		unlockSign, err = harmonic.Nineora().TokenUnlock(innerTx, &nineapi.TxUnlockParas{
			Account: localJob.Payer,
			Amount:  localJob.Amount,
			Biz:     pay.BizCode,
			Link:    bizID,
			Ex:      localJob.Ex,
		})
		if err != nil {
			hlog.Err("hyper.payment.nine.Advance:do unlock failed", zap.Error(err))
			return err
		}
		transSign, err = harmonic.Nineora().TokenTransfer(innerTx, &nineapi.TxTransferParas{
			Sender:     localJob.Payer,
			Recipient:  localJob.Payee,
			Amount:     localJob.Amount,
			LockAmount: 0,
			Biz:        pay.BizCode,
			Link:       bizID,
			Ex:         localJob.Ex,
		})
		if err != nil {
			hlog.Err("hyper.payment.nine.Advance:do transfer failed", zap.Error(err))
			return err
		}
		return nil
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

func (e *NineExecutor) Cancel(ctx context.Context, pay *Payment, job *Job) (synced bool, err error) {
	//todo
	hlog.TraceFix("no support to cancel job right now!",
		ctx,
		errors.New("no support to cancel job right now"),
		zap.Any("pay", pay),
		zap.Any("job", job))
	return true, nil
}

func (e *NineExecutor) Timeout(_ context.Context, _ *Payment, _ *Job) (synced bool, err error) {
	return true, nil
}

func (e *NineExecutor) OnPrepared(ctx context.Context, job *Job) error {
	fmt.Println("OnPrepared job: ", hjson.MustToString(job))
	return nil
}

func (e *NineExecutor) OnCompleted(ctx context.Context, job *Job) error {
	fmt.Println("OnCompleted job: ", hjson.MustToString(job))
	return nil
}

func (e *NineExecutor) OnTimeout(ctx context.Context, job *Job) error {
	fmt.Println("OnTimeout job: ", hjson.MustToString(job))
	return nil
}

func (e *NineExecutor) OnCanceled(ctx context.Context, job *Job) error {
	fmt.Println("OnCanceled job: ", hjson.MustToString(job))
	return nil
}
