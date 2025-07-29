package payment

import (
	"context"
	"errors"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/hyperplt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func onPayJobPrepared(ctx context.Context, jobM *JobM) (err error) {
	if jobM == nil {
		hlog.TraceErr("hyper.payment.onPayJobPrepared: jobM == nil", ctx, errors.New("jobM is nil"))
		return nil
	}
	if hlog.IsElapseDetail() {
		defer hlog.ElapseWithCtx(ctx, "hyper.payment.onPayJobPrepared",
			hlog.F(zap.String("job_id", jobM.ID)),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err)}
				}
				return nil
			})()
	}
	tx := hyperplt.Tx(ctx)
	payM, err := hdb.Get[PayM](tx, "id = ?", jobM.PaymentID)
	if err != nil {
		return err
	}
	if payM == nil {
		hlog.TraceFix("hyper.payment.onPayJobPrepared",
			ctx, errors.New("no such payment"),
			zap.String("job_id", jobM.ID),
			zap.Uint64("payment_id", jobM.PaymentID))
		return nil
	}

	lstPayAutoID := int64(0)
	preparedJobCount := 0
	err = hdb.Iter[JobM](func() *gorm.DB {
		return tx.Where("payment_id = ? AND auto_id > ?", jobM.PaymentID, lstPayAutoID)
	}, func(m *JobM) error {
		switch m.Status {
		case JobPrepared, JobCompleted:
			preparedJobCount++
		default:
		}
		lstPayAutoID = m.AutoID
		return nil
	})
	if err != nil {
		return errors.New("load job failed" + err.Error())
	}

	if jobM.Timeout > 0 {
		ttListenJobTimeout(ctx, jobM.To())
	}

	if payM.JobCount == preparedJobCount {
		return doAdvToPrepared(ctx, jobM.PaymentID)
	}

	return nil
}

func onPayJobTimeout(ctx context.Context, jobM *JobM) error {
	return doAdvToTimeout(ctx, jobM.PaymentID)
}

func onPayJobCanceled(ctx context.Context, jobM *JobM) error {
	return doAdvToCanceled(ctx, jobM.PaymentID)
}

func onPayJobCompleted(ctx context.Context, jobM *JobM) (err error) {
	if jobM == nil {
		hlog.TraceErr("hyper.payment.onPayJobCompleted: jobM == nil", ctx, errors.New("jobM is nil"))
		return nil
	}
	if hlog.IsElapseDetail() {
		defer hlog.ElapseWithCtx(ctx, "hyper.payment.onPayJobCompleted",
			hlog.F(zap.String("job_id", jobM.ID)),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err)}
				}
				return nil
			})()
	}
	tx := hyperplt.Tx(ctx)
	payM, err := hdb.Get[PayM](tx, "id = ?", jobM.PaymentID)
	if err != nil {
		return err
	}
	if payM == nil {
		hlog.TraceFix("hyper.payment.onPayJobCompleted",
			ctx, errors.New("no such payment"),
			zap.String("job_id", jobM.ID),
			zap.Uint64("payment_id", jobM.PaymentID))
		return nil
	}

	lstPayAutoID := int64(0)
	completedJobCount := 0
	err = hdb.Iter[JobM](func() *gorm.DB {
		return tx.Where("payment_id = ? AND auto_id > ?", jobM.PaymentID, lstPayAutoID)
	}, func(m *JobM) error {
		switch m.Status {
		case JobCompleted:
			completedJobCount++
		default:
		}
		lstPayAutoID = m.AutoID
		return nil
	})
	if err != nil {
		return errors.New("load job failed" + err.Error())
	}
	if payM.JobCount == completedJobCount {
		return doAdvToCompleted(ctx, jobM.PaymentID)
	}

	return nil
}
