package payment

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
	paymentID := cast.ToUint64(data[0].Get("payment_id"))
	err := doSetStatus(ctx, paymentID, src, Completed)
	if err != nil {
		return 0, err
	}
	//todo do completed
	return Completed, nil
}

func doAdvToCompleted(ctx context.Context, paymentID ID) error {
	return doAdvance(ctx, paymentID, CompleteEvent)
}
