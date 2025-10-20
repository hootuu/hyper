package payment

import (
	"context"
	"errors"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/hyperplt"
	"go.uber.org/zap"
)

func onJobAlter(ctx context.Context, jobID JobID, src JobStatus, dst JobStatus) (err error) {
	if hlog.IsElapseFunction() {
		defer hlog.ElapseWithCtx(ctx, "hyper.payment.onJobAlter",
			hlog.F(zap.String("job_id", jobID), zap.Int("src", int(src)), zap.Int("dst", int(dst))),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err)}
				}
				return nil
			},
		)()
	}
	jobM, err := hdb.Get[JobM](hyperplt.Tx(ctx), "id=? AND status=?", jobID, dst)
	if err != nil {
		return err
	}
	if jobM == nil {
		hlog.TraceFix("hyper.payment.onJobAlter",
			ctx, errors.New("no such job"),
			zap.String("id", jobID),
			zap.Int("src", int(src)),
			zap.Int("dst", int(dst)))
		return nil
	}
	switch dst {
	case JobTimeout:
		err = onPayJobTimeout(ctx, jobM)
	case JobPrepared:
		err = onPayJobPrepared(ctx, jobM)
	case JobCanceled:
		err = onPayJobCanceled(ctx, jobM)
	case JobCompleted:
		err = onPayJobCompleted(ctx, jobM)
	default:
	}
	if err != nil {
		return err
	}
	return nil
}

func UpdateJobPayType(ctx context.Context, paymentID uint64, payType string) (err error) {
	payM, err := hdb.MustGet[JobM](hyperplt.DB(), "payment_id = ?", paymentID)
	if err != nil {
		return err
	}
	ext := *hjson.MustFromBytes[dict.Dict](payM.Ctx)
	ext["third_code"] = payType
	return hdb.Update[JobM](hyperplt.Tx(ctx), map[string]any{"ctx": hjson.MustToBytes(ext)}, "auto_id = ?", payM.AutoID)
}
