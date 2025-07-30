package hiorder

import (
	"context"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/hfsm"
	"github.com/spf13/cast"
)

func (e *Engine[T]) advToConsensus(
	ctx context.Context,
	src hfsm.State,
	_ hfsm.Event,
	data ...dict.Dict,
) (hfsm.State, error) {
	orderID := cast.ToUint64(data[0].Get("order_id").String())
	err := e.doSetStatus(ctx, orderID, Consensus, src, nil)
	if err != nil {
		return 0, err
	}
	return Consensus, nil
}

func (e *Engine[T]) doAdvToConsensus(
	ctx context.Context,
) (err error) {
	return e.doAdvance(ctx, ConsenseEvent, nil)
}
