package payment

import (
	"context"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/hfsm"
	"github.com/spf13/cast"
)

func advToTimeout(
	ctx context.Context,
	src hfsm.State,
	_ hfsm.Event,
	data ...dict.Dict,
) (hfsm.State, error) {
	paymentID := cast.ToUint64(data[0].Get("payment_id").String())
	err := doSetStatus(ctx, paymentID, Timeout, src)
	if err != nil {
		return 0, err
	}
	return Timeout, nil
}

func doAdvToTimeout(ctx context.Context, paymentID ID) error {
	return doAdvance(ctx, paymentID, TimeoutEvent)
}
