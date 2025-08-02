package payment

import (
	"context"
	"errors"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/hfsm"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/hyperplt"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func advToPrepared(
	ctx context.Context,
	src hfsm.State,
	_ hfsm.Event,
	data ...dict.Dict,
) (hfsm.State, error) {
	paymentID := cast.ToUint64(data[0].Get("payment_id").String())
	err := doSetStatus(ctx, paymentID, Prepared, src)
	if err != nil {
		return 0, err
	}
	return Prepared, nil
}

func doAdvToPrepared(ctx context.Context, paymentID ID) error {
	return doAdvance(ctx, paymentID, PrepareEvent)
}

func doPrepare(ctx context.Context, payID ID) (err error) {
	if hlog.IsElapseComponent() {
		defer hlog.ElapseWithCtx(ctx,
			"hyper.payment.doPrepare",
			hlog.F(zap.Uint64("payID", payID)),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err)}
				}
				return nil
			},
		)()
	}

	tx := hyperplt.Tx(ctx)
	payM, err := hdb.MustGet[PayM](tx, "id = ?", payID)
	if err != nil {
		return errors.New("payment not exists: " + err.Error())
	}
	pay := payM.To()

	lstPayAutoID := int64(0)
	var jobArrM []*JobM
	err = hdb.Iter[JobM](func() *gorm.DB {
		return tx.Order("auto_id asc").Where("payment_id = ? AND auto_id > ?", payID, lstPayAutoID)
	}, func(m *JobM) error {
		jobArrM = append(jobArrM, m)
		lstPayAutoID = m.AutoID
		return nil
	})
	if err != nil {
		return errors.New("load job failed" + err.Error())
	}
	if len(jobArrM) == 0 {
		return errors.New("load job failed, no job found")
	}
	err = hdb.Tx(tx, func(tx *gorm.DB) error {
		preparedCount := 0
		innerCtx := hdb.TxCtx(tx, ctx)
		for _, job := range jobArrM {
			executor, innerErr := MustGetJobExecutor(job.Channel)
			if innerErr != nil {
				return innerErr
			}
			done, innerErr := executor.Prepare(innerCtx, pay, job.To())
			if innerErr != nil {
				return innerErr
			}
			if done {
				preparedCount++
			}
		}

		if preparedCount > 0 {
			if preparedCount >= payM.JobCount {
				err = doAdvance(innerCtx, payID, PrepareEvent)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		hlog.TraceErr("hyper.payment.doPrepare", ctx, err)
		return err
	}

	return nil
}
