package hiorder

import (
	"context"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/hfsm"
	"github.com/spf13/cast"
)

func (e *Engine[T]) advToCanceled(
	ctx context.Context,
	src hfsm.State,
	_ hfsm.Event,
	data ...dict.Dict,
) (hfsm.State, error) {
	orderID := cast.ToUint64(data[0].Get("order_id").String())
	err := e.doSetStatus(ctx, orderID, Canceled, src, nil)
	if err != nil {
		return 0, err
	}
	return Canceled, nil
}

func (e *Engine[T]) doAdvToCanceled(
	ctx context.Context,
	orderID ID,
) (err error) {
	return e.doAdvance(ctx, orderID, CancelEvent, nil)
}
