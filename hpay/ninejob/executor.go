package ninejob

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/hpay/payment"
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

type Executor struct{}

func NewExecutor() *Executor {
	return &Executor{}
}

func (e *Executor) GetChannel() payment.Channel {
	return NineChannel
}

func (e *Executor) Prepare(ctx context.Context, pay *payment.Payment, job *payment.Job) (err error) {
	if pay == nil {
		return fmt.Errorf("assert: pay is not nil")
	}
	if job == nil {
		return fmt.Errorf("assert: job is not nil")
	}
	if pay.ID != job.PaymentID {
		return fmt.Errorf("assert: payment from job %d is not %d", job.PaymentID, pay.ID)
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

	localJob, err := JobFromCtx(job.Ctx)
	if err != nil {
		return err
	}

	bizCode, _, _ := pay.Biz.ToCodeID()
	_, err = harmonic.Nineora().TokenLock(ctx, &nineapi.TxLockParas{
		Account: localJob.Payer,
		Amount:  localJob.Amount,
		Biz:     bizCode,
		Link:    pay.Biz,
		Ex:      localJob.Ex,
	})
	if err != nil {
		return err
	}

	fmt.Println("Prepare job: ", hjson.MustToString(job)) //todo

	return nil
}

func (e *Executor) Advance(ctx context.Context, pay *payment.Payment, job *payment.Job) (err error) {
	if pay == nil {
		return fmt.Errorf("assert: pay is not nil")
	}
	if job == nil {
		return fmt.Errorf("assert: job is not nil")
	}
	if pay.ID != job.PaymentID {
		return fmt.Errorf("assert: payment from job %d is not %d", job.PaymentID, pay.ID)
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

	localJob, err := JobFromCtx(job.Ctx)
	if err != nil {
		return err
	}

	bizCode, _, _ := pay.Biz.ToCodeID()

	err = hdb.Tx(hyperplt.Tx(ctx), func(tx *gorm.DB) error {
		innerTx := hdb.TxCtx(tx, ctx)
		unlockSign, err = harmonic.Nineora().TokenUnlock(innerTx, &nineapi.TxUnlockParas{
			Account: localJob.Payer,
			Amount:  localJob.Amount,
			Biz:     bizCode,
			Link:    pay.Biz,
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
			Biz:        bizCode,
			Link:       pay.Biz,
			Ex:         localJob.Ex,
		})
		if err != nil {
			hlog.Err("hyper.payment.nine.Advance:do transfer failed", zap.Error(err))
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (e *Executor) Cancel(ctx context.Context, job *payment.Job) error {
	fmt.Println("Cancel job: ", hjson.MustToString(job))
	return nil
}

func (e *Executor) OnPrepared(ctx context.Context, job *payment.Job) error {
	fmt.Println("OnPrepared job: ", hjson.MustToString(job))
	return nil
}

func (e *Executor) OnCompleted(ctx context.Context, job *payment.Job) error {
	fmt.Println("OnCompleted job: ", hjson.MustToString(job))
	return nil
}

func (e *Executor) OnTimeout(ctx context.Context, job *payment.Job) error {
	fmt.Println("OnTimeout job: ", hjson.MustToString(job))
	return nil
}

func (e *Executor) OnCanceled(ctx context.Context, job *payment.Job) error {
	fmt.Println("OnCanceled job: ", hjson.MustToString(job))
	return nil
}
