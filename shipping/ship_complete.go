package shipping

import (
	"context"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/hfsm"
	"github.com/spf13/cast"
)

func advToCompleted(
	ctx context.Context,
	src hfsm.State,
	_ hfsm.Event,
	data ...dict.Dict,
) (hfsm.State, error) {
	shippingID := cast.ToUint64(data[0].Get("shipping_id").String())
	timeoutCompleted := cast.ToBool(data[0].Get("timeout_completed").String())
	var innerData dict.Dict = nil
	if timeoutCompleted {
		innerData = dict.NewDict().Set("timeout_completed", true)
	}
	err := doSetStatus(ctx, shippingID, Completed, src, innerData)
	if err != nil {
		return 0, err
	}
	return Completed, nil
}

func doAdvToCompleted(ctx context.Context, shippingID ID) error {
	return doAdvance(ctx, shippingID, CompleteEvent, nil)
}

func doAdvToTimeoutCompleted(ctx context.Context, shippingID ID) error {
	return doAdvance(ctx, shippingID, CompleteEvent,
		func(shipM *ShipM, mustMut map[string]any) {
			mustMut["timeout_completed"] = true
		},
	)
}
