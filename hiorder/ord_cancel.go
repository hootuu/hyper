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
	var mut dict.Dict
	if len(data) > 1 {
		mut = data[1]
	}
	err := e.doSetStatus(ctx, orderID, Canceled, src, mut)
	if err != nil {
		return 0, err
	}
	return Canceled, nil
}

func (e *Engine[T]) doAdvToCanceled(
	ctx context.Context,
	mutSet func(ordM *OrderM, mustMut map[string]any),
) (err error) {
	return e.doAdvance(ctx, CancelEvent, mutSet)
}
