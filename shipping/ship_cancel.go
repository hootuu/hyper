package shipping

import (
	"context"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/hfsm"
	"github.com/spf13/cast"
)

func advToCanceled(
	ctx context.Context,
	src hfsm.State,
	_ hfsm.Event,
	data ...dict.Dict,
) (hfsm.State, error) {
	shippingID := cast.ToUint64(data[0].Get("shipping_id").String())
	err := doSetStatus(ctx, shippingID, Canceled, src, nil)
	if err != nil {
		return 0, err
	}
	return Canceled, nil
}

func doAdvToCanceled(ctx context.Context, shippingID ID) error {
	return doAdvance(ctx, shippingID, CompleteEvent, nil)
}
