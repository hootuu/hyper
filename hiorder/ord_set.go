package hiorder

import (
	"context"
	"errors"
	"fmt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/idx"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/ex"
	"github.com/hootuu/hyper/hyperplt"
	"github.com/hootuu/hyper/payment"
	"github.com/hootuu/hyper/shipping"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

func (e *Engine[T]) SetShipping(ctx context.Context, ordID ID, shippingID shipping.ID) (err error) {
	if hlog.IsElapseFunction() {
		defer hlog.ElapseWithCtx(ctx, "hyper.order.eng.SetShipping",
			hlog.F(zap.Uint64("ordID", ordID), zap.Uint64("shippingID", shippingID)),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err)}
				}
				return nil
			},
		)()
	}
	mut := map[string]any{
		"shipping_id": shippingID,
	}
	err = hdb.Update[OrderM](hyperplt.Tx(ctx), mut, "id = ?", ordID)
	if err != nil {
		return errors.New("set shipping id failed: " + err.Error())
	}
	return nil
}

type SetPaymentParas struct {
	Idem     string              `json:"idem"`
	OrderID  ID                  `json:"order_id"`
	Payments []payment.JobDefine `json:"payments"`
	Timeout  time.Duration       `json:"timeout"`
	Ex       *ex.Ex              `json:"ex"`
}

func (e *Engine[T]) SetPayment(ctx context.Context, paras *SetPaymentParas) (paymentID payment.ID, err error) {
	if hlog.IsElapseComponent() {
		defer hlog.ElapseWithCtx(ctx, "hyper.order.SetPayment",
			hlog.F(zap.Uint64("ordId", paras.OrderID)),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err), zap.Any("payment", paras.Payments)}
				}
				return nil
			},
		)()
	}
	if len(paras.Payments) == 0 {
		return 0, fmt.Errorf("require payment")
	}
	if e.ord.PaymentID != 0 {
		//TODO fix with trice
		hlog.Info("ordM.PaymentID!=0", hlog.TraceInfo(ctx), zap.Error(err))
		return e.ord.PaymentID, nil
	}
	if err := hyperplt.Idem().MustCheck(paras.Idem); err != nil {
		return 0, err
	}
	err = hdb.Tx(hyperplt.Tx(ctx), func(tx *gorm.DB) error {
		innerCtx := hdb.TxCtx(tx, ctx)
		paymentID, err = payment.Create(innerCtx, &payment.CreateParas{
			Idem:    idx.New(),
			Payer:   e.ord.Payer,
			Payee:   e.ord.Payee,
			BizCode: e.ord.Code,
			BizID:   cast.ToString(paras.OrderID),
			Amount:  e.ord.Amount,
			Ex:      paras.Ex,
			Jobs:    paras.Payments,
			Timeout: paras.Timeout,
		})
		if err != nil {
			hlog.Err("Create err", hlog.TraceInfo(ctx), zap.Error(err))
			return errors.New("Create error: " + err.Error())
		}
		err = payment.Prepare(innerCtx, paymentID)
		if err != nil {
			hlog.Err("payment.Prepare err", hlog.TraceInfo(ctx), zap.Error(err))
			return errors.New("payment prepared: " + err.Error())
		}
		ordMut := map[string]any{
			"payment_id": paymentID,
		}
		err = hdb.Update[OrderM](tx, ordMut, "id = ?", paras.OrderID)
		if err != nil {
			hlog.Err("Update[OrderM] err", hlog.TraceInfo(ctx), zap.Error(err))
			return errors.New("Update Order Failed: " + err.Error())
		}
		return nil
	})
	if err != nil {
		hlog.Err("Tx err", hlog.TraceInfo(ctx), zap.Error(err))
		return 0, errors.New("Tx Failed: " + err.Error())
	}
	e.ord.PaymentID = paymentID
	return paymentID, nil
}
