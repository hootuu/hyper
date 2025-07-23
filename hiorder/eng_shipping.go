package hiorder

import (
	"context"
	"errors"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/hshipping/shipping"
	"github.com/hootuu/hyper/hyperplt"
	"go.uber.org/zap"
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
