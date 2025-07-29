package shipping

import (
	"context"
	"errors"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/hfsm"
	"github.com/spf13/cast"
)

func advToSubmitted(
	ctx context.Context,
	src hfsm.State,
	_ hfsm.Event,
	data ...dict.Dict,
) (hfsm.State, error) {
	shippingID := cast.ToUint64(data[0].Get("shipping_id").String())
	err := doSetStatus(ctx, shippingID, Submitted, src, nil)
	if err != nil {
		return 0, err
	}
	return Submitted, nil
}

func doAdvToSubmitted(
	ctx context.Context,
	shippingID ID,
	courierCode string,
	trackingNo string,
) (err error) {
	if courierCode == "" {
		return errors.New("courier_code is required")
	}
	if trackingNo == "" {
		return errors.New("tracking_no is required")
	}
	return doAdvance(ctx, shippingID, SubmitEvent,
		func(shipM *ShipM, mustMut map[string]any) {
			mustMut["courier_code"] = courierCode
			mustMut["tracking_no"] = trackingNo
		},
	)
}
