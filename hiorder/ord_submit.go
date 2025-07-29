package hiorder

import (
	"context"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/hfsm"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/hyperplt"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (e *Engine[T]) advToInitial(
	ctx context.Context,
	src hfsm.State,
	_ hfsm.Event,
	data ...dict.Dict,
) (hfsm.State, error) {
	orderID := cast.ToUint64(data[0].Get("order_id").String())
	err := e.doSetStatus(ctx, orderID, Initial, src, nil)
	if err != nil {
		return 0, err
	}
	return Initial, nil
}

func (e *Engine[T]) doAdvToInitial(
	ctx context.Context,
	orderID ID,
) (err error) {
	ordM := e.ord.toModel()

	tx := hyperplt.Tx(ctx)
	err = hdb.Tx(tx, func(tx *gorm.DB) error {
		innerCtx := hdb.TxCtx(tx)
		err = e.deal.Before(innerCtx, e.ord.Status, Initial)
		if err != nil {
			hlog.Err("hyper.order.Submit: dealer.Before", zap.Error(err))
			return err
		}
		ordM.Status = Initial
		err = hdb.Create[OrderM](tx, ordM)
		if err != nil {
			hlog.Err("hyper.order.Submit: db.Create", zap.Error(err))
			return err
		}
		err = e.deal.After(innerCtx, e.ord.Status, Initial)
		if err != nil {
			hlog.Err("hyper.order.Submit: dealer.After", zap.Error(err))
			return err
		}
		return nil
	})
	if err != nil {
		hlog.Err("hyper.order.Submit: Tx", zap.Error(err))
		return err
	}

	return nil
}
