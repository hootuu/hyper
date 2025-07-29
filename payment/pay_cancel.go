package payment

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
	paymentID := cast.ToUint64(data[0].Get("payment_id").String())
	err := doSetStatus(ctx, paymentID, src, Canceled)
	if err != nil {
		return 0, err
	}
	//canceled
	return Canceled, nil
}

func doAdvToCanceled(ctx context.Context, paymentID ID) error {
	return doAdvance(ctx, paymentID, CancelEvent)
}
