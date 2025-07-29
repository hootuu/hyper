package shipping

import (
	"context"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/hfsm"
	"github.com/spf13/cast"
)

func advToFailed(
	ctx context.Context,
	src hfsm.State,
	_ hfsm.Event,
	data ...dict.Dict,
) (hfsm.State, error) {
	shippingID := cast.ToUint64(data[0].Get("shipping_id").String())
	err := doSetStatus(ctx, shippingID, Failed, src, nil)
	if err != nil {
		return 0, err
	}
	return Failed, nil
}

func doAdvToFailed(ctx context.Context, shippingID ID) error {
	return doAdvance(ctx, shippingID, FailEvent, nil)
}
