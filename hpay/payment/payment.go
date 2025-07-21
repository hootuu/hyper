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
)

type CreateParas struct {
	Payer   collar.Link `json:"payer"`
	Payee   collar.Link `json:"payee"`
	BizLink collar.Link `json:"biz"`
	Amount  uint64      `json:"amount"`
	Ex      *ex.Ex      `json:"ex"`
	Jobs    []JobDefine `json:"jobs"`
}

func (paras *CreateParas) Validate() error {
	if paras.BizLink == "" {
		return errors.New("biz_link is required")
	}
	if paras.Payer == "" {
		return errors.New("payer is required")
	}
	if paras.Payee == "" {
		return errors.New("payee is required")
	}
	if paras.Amount == 0 {
		return errors.New("amount is required")
	}
	if len(paras.Jobs) == 0 {
		return errors.New("jobs is required")
	}
	return nil
}

func Create(ctx context.Context, paras *CreateParas) (id ID, err error) {
	if paras == nil {
		return 0, fmt.Errorf("hpay.Create: assert paras != nil")
	}
	if err := paras.Validate(); err != nil {
		return 0, err
	}
	if hlog.IsElapseComponent() {
		defer hlog.ElapseWithCtx(ctx, "hyper.payment.Create",
			hlog.F(zap.String("biz", paras.BizLink.Str())),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{
						zap.String("biz", paras.BizLink.Display()),
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
	tx := hyperplt.Tx(ctx)
	bBizExist, err := hdb.Exist[PayM](tx, "biz_link = ?", paras.BizLink)
	if err != nil {
		return 0, err
	}
	if bBizExist {
		return 0, fmt.Errorf("hpay.Create: biz_link exist[%s]", paras.BizLink.Display())
	}

	payM := &PayM{
		Template: hdb.TemplateFromEx(paras.Ex),
		ID:       NxtPaymentID(),
		BizLink:  paras.BizLink,
		Payer:    paras.Payer,
		Payee:    paras.Payee,
		Amount:   paras.Amount,
		Status:   Initialized,
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
		seq := uint32(i + 1)
		jobM := &JobM{
			ID:         BuildJobID(payM.ID, seq),
			Channel:    job.GetChannel(),
			PaymentID:  payM.ID,
			PaymentSeq: seq,
			Status:     JobInitialized,
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

func Prepare(ctx context.Context, payID ID) (err error) {
	if hlog.IsElapseComponent() {
		defer hlog.ElapseWithCtx(ctx, "hyper.payment.Prepare",
			hlog.F(zap.Uint64("payID", payID)),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err)}
				}
				return nil
			},
		)()
	}
	var jobArrM []*JobM
	tx := hyperplt.Tx(ctx)
	payM, err := hdb.MustGet[PayM](tx, "id = ?", payID)
	if err != nil {
		return errors.New("payment not exists: " + err.Error())
	}
	pay := payM.To()
	lstPayAutoID := int64(0)
	err = hdb.Iter[JobM](func() *gorm.DB {
		return tx.Where("payment_id = ? AND auto_id > ?", payID, lstPayAutoID)
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
		innerCtx := hdb.TxCtx(tx, ctx)
		for _, job := range jobArrM {
			executor, err := MustGetJobExecutor(job.Channel)
			if err != nil {
				return err
			}
			err = executor.Prepare(innerCtx, pay, job.To())
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return errors.New("job executor" + err.Error())
	}

	return nil
}

func Advance(ctx context.Context, payID ID) (err error) {
	if hlog.IsElapseComponent() {
		defer hlog.ElapseWithCtx(ctx, "hyper.payment.Advance",
			hlog.F(zap.Uint64("payID", payID)),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err)}
				}
				return nil
			},
		)()
	}
	var jobArrM []*JobM
	tx := hyperplt.Tx(ctx)
	payM, err := hdb.MustGet[PayM](tx, "id = ?", payID)
	if err != nil {
		return errors.New("payment not exists: " + err.Error())
	}
	pay := payM.To()
	lstPayAutoID := int64(0)
	err = hdb.Iter[JobM](func() *gorm.DB {
		return tx.Where("payment_id = ? AND auto_id > ?", payID, lstPayAutoID)
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
		innerCtx := hdb.TxCtx(tx, ctx)
		for _, job := range jobArrM {
			executor, err := MustGetJobExecutor(job.Channel)
			if err != nil {
				return err
			}
			err = executor.Advance(innerCtx, pay, job.To())
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return errors.New("job executor" + err.Error())
	}

	return nil
}
