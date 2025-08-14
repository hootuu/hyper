package payment

import (
	"context"
	"errors"
	"fmt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyle/hypes/ex"
	"github.com/hootuu/hyper/hyperplt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type CreateParas struct {
	Idem    string        `json:"idem"`
	Payer   collar.Link   `json:"payer"`
	Payee   collar.Link   `json:"payee"`
	BizCode string        `json:"biz_code"`
	BizID   string        `json:"biz_id"`
	Amount  uint64        `json:"amount"`
	Ex      *ex.Ex        `json:"ex"`
	Jobs    []JobDefine   `json:"jobs"`
	Timeout time.Duration `json:"timeout"`
}

func (paras *CreateParas) Validate() error {
	if paras.Idem == "" {
		return errors.New("idem is required")
	}
	if paras.BizCode == "" {
		return errors.New("biz_code is required")
	}
	if paras.BizID == "" {
		return errors.New("biz_id is required")
	}
	if paras.Payer == "" {
		return errors.New("payer is required")
	}
	if paras.Payee == "" {
		return errors.New("payee is required")
	}
	//if paras.Amount == 0 {
	//	return errors.New("amount is required")
	//}
	if len(paras.Jobs) == 0 {
		return errors.New("jobs is required")
	}
	if paras.Timeout == 0 {
		paras.Timeout = 15 * time.Minute
	}
	return nil
}

func Create(ctx context.Context, paras *CreateParas) (id ID, err error) {
	InitIfNeeded()
	if paras == nil {
		return 0, fmt.Errorf("hpay.Create: assert paras != nil")
	}
	if err := paras.Validate(); err != nil {
		return 0, err
	}
	if hlog.IsElapseComponent() {
		defer hlog.ElapseWithCtx(ctx, "hyper.payment.Create",
			hlog.F(
				zap.String("biz_code", paras.BizCode),
				zap.String("biz", paras.BizID),
			),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{
						zap.String("biz_code", paras.BizCode),
						zap.String("biz", paras.BizID),
						zap.String("payer", paras.Payer.Display()),
						zap.String("payee", paras.Payee.Display()),
						zap.Uint64("amount", paras.Amount),
						zap.Error(err),
					}
				}
				return []zap.Field{zap.Uint64("id", id)}
			},
		)()
	}
	if err := hyperplt.Idem().MustCheck(paras.Idem); err != nil {
		return 0, err
	}
	tx := hyperplt.Tx(ctx)
	payM := &PayM{
		Template: hdb.TemplateFromEx(paras.Ex),
		ID:       NxtPaymentID(),
		BizCode:  paras.BizCode,
		BizID:    paras.BizID,
		Payer:    paras.Payer,
		Payee:    paras.Payee,
		Amount:   paras.Amount,
		Status:   Initialized,
		Timeout:  paras.Timeout,
		JobCount: len(paras.Jobs),
	}
	totalAmountByJob := uint64(0)
	var jobArrM []*JobM
	for i, job := range paras.Jobs {
		if err := job.Validate(); err != nil {
			return 0, err
		}
		totalAmountByJob += job.GetAmount()
		if totalAmountByJob > paras.Amount {
			return 0, errors.New("hpay.Create: totalAmountByJobs exceeds amount")
		}
		seq := i + 1
		jobM := &JobM{
			ID:         BuildJobID(payM.ID, seq),
			Channel:    job.GetChannel(),
			PaymentID:  payM.ID,
			PaymentSeq: seq,
			Status:     JobInitialized,
			Timeout:    paras.Timeout,
			Ctx:        hjson.MustToBytes(job.GetCtx()),
		}
		jobArrM = append(jobArrM, jobM)
	}
	if totalAmountByJob != paras.Amount {
		return 0, errors.New("total amount calced by jobs != amount")
	}
	err = hdb.Tx(tx, func(tx *gorm.DB) error {
		err = hdb.Create[PayM](tx, payM)
		if err != nil {
			return err
		}
		err = hdb.MultiCreate[JobM](tx, jobArrM)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, errors.New("db.Tx: failed" + err.Error())
	}

	return payM.ID, nil
}
