package payment

import (
	"context"
	"errors"
	"fmt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/hfsm"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/hyperplt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var gPaymentFSM *hfsm.Machine
var gJobFSM *hfsm.Machine

func init() {
	gPaymentFSM = hfsm.NewMachine().
		AddTransition(Initialized, PrepareEvent, advToPrepared).
		AddTransition(Initialized, CancelEvent, advToCanceled).
		AddTransition(Initialized, TimeoutEvent, advToTimeout).
		AddTransition(Prepared, TimeoutEvent, advToTimeout).
		AddTransition(Prepared, CancelEvent, advToCanceled).
		AddTransition(Prepared, ExecuteEvent, advToCompleted).
		AddTransition(Prepared, CompleteEvent, advToCompleted).
		AddTransition(Executing, CompleteEvent, advToCompleted).
		AddTransition(Executing, TimeoutEvent, advToTimeout)

	gJobFSM = hfsm.NewMachine().
		AddTransition(JobInitialized, JobPrepareEvent, advJobToPrepared).
		AddTransition(JobInitialized, JobCancelEvent, advJobToCanceled).
		AddTransition(JobPrepared, JobTimeoutEvent, advJobToTimeout).
		AddTransition(JobPrepared, JobCompleteEvent, advJobToCompleted).
		AddTransition(JobPrepared, JobCancelEvent, advJobToCanceled)
}

func doAdvance(ctx context.Context, paymentID ID, event hfsm.Event) (err error) {
	if hlog.IsElapseDetail() {
		defer hlog.ElapseWithCtx(ctx, "hyper.payment.doAdvance",
			hlog.F(zap.Uint64("paymentID", paymentID), zap.Int("event", int(event))),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err)}
				}
				return nil
			},
		)()
	}
	var payM *PayM
	payM, err = hdb.Get[PayM](hyperplt.Tx(ctx), "id = ?", paymentID)
	if err != nil {
		return err
	}
	if payM == nil {
		hlog.TraceFix("payment.doAdvance: no such payment", ctx, errors.New("no such payment"),
			zap.Uint64("payment_id", paymentID))
		return nil
	}
	_, err = gPaymentFSM.Handle(ctx, payM.Status, event, dict.NewDict().Set("payment_id", paymentID))
	if err != nil {
		return err
	}
	return nil
}

func doJobAdvance(ctx context.Context, jobID JobID, event hfsm.Event, data dict.Dict) error {
	jobM, err := hdb.Get[JobM](hyperplt.Tx(ctx), "id = ?", jobID)
	if err != nil {
		return err
	}
	if jobM == nil {
		hlog.TraceFix("payment.doJobAdvance: no such job", ctx, errors.New("no such job"),
			zap.String("job_id", jobID))
		return fmt.Errorf("no such job: %s", jobID)
	}
	if data == nil {
		data = dict.NewDict()
	}
	data.Set("job_id", jobID)
	_, err = gJobFSM.Handle(ctx, jobM.Status, event, data)
	if err != nil {
		return err
	}
	return nil
}

func doSetStatus(ctx context.Context, id ID, targetStatus Status, srcStatus Status) (err error) {
	if hlog.IsElapseComponent() {
		defer hlog.ElapseWithCtx(ctx, "hyper.payment.doSetStatus",
			hlog.F(zap.Uint64("id", id), zap.Int("target", int(targetStatus)),
				zap.Int("src", int(srcStatus))),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err)}
				}
				return nil
			},
		)()
	}
	tx := hyperplt.Tx(ctx)
	payM, err := hdb.MustGet[PayM](tx, "id = ?", id)
	if err != nil {
		return err
	}

	mut := map[string]interface{}{
		"status": targetStatus,
	}
	switch targetStatus {
	case Prepared:
		mut["prepared_time"] = gorm.Expr("CURRENT_TIMESTAMP")
	case Executing:
		mut["executing_time"] = gorm.Expr("CURRENT_TIMESTAMP")
	case Timeout:
		mut["timeout_time"] = gorm.Expr("CURRENT_TIMESTAMP")
	case Canceled:
		mut["canceled_time"] = gorm.Expr("CURRENT_TIMESTAMP")
	case Completed:
		mut["completed_time"] = gorm.Expr("CURRENT_TIMESTAMP")
	default:
	}
	rows, err := hdb.UpdateX[PayM](tx, mut, "id = ? AND status = ?", id, srcStatus)
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("payment[id=%d, status=%d] not exist", id, srcStatus)
	}
	mqPublishPaymentAlter(ctx, &AlterPayload{
		PaymentID: payM.ID,
		BizCode:   payM.BizCode,
		BizID:     payM.BizID,
		Src:       srcStatus,
		Dst:       targetStatus,
	})
	return nil
}

func doSetJobStatus(
	ctx context.Context,
	jobID JobID,
	targetStatus Status,
	srcStatus Status,
	data dict.Dict,
) (err error) {
	if hlog.IsElapseComponent() {
		defer hlog.ElapseWithCtx(ctx, "hyper.payment.doSetJobStatus",
			hlog.F(zap.String("jobID", jobID), zap.Int("target", int(targetStatus)),
				zap.Int("src", int(srcStatus))),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err)}
				}
				return nil
			},
		)()
	}
	tx := hyperplt.Tx(ctx)
	jobM, err := hdb.MustGet[JobM](tx, "id = ?", jobID)
	if err != nil {
		return err
	}
	mut := data
	if mut == nil {
		mut = dict.NewDict()
	}
	mut["status"] = targetStatus
	switch targetStatus {
	case JobPrepared:
		mut["prepared_time"] = gorm.Expr("CURRENT_TIMESTAMP")
	case JobTimeout:
		mut["timeout_time"] = gorm.Expr("CURRENT_TIMESTAMP")
	case JobCanceled:
		mut["canceled_time"] = gorm.Expr("CURRENT_TIMESTAMP")
	case JobCompleted:
		mut["completed_time"] = gorm.Expr("CURRENT_TIMESTAMP")
	default:
	}

	rows, err := hdb.UpdateX[JobM](tx, mut, "id = ? AND status = ?", jobID, srcStatus)
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("job[id=%s, status=%d] not exist", jobID, srcStatus)
	}

	mqPublishJobAlter(ctx, &JobAlterPayload{
		JobID: jobM.ID,
		Src:   srcStatus,
		Dst:   targetStatus,
	})
	return nil
}
