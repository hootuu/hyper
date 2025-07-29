package hiorder

import (
	"context"
	"fmt"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/payment"
	"github.com/spf13/cast"
)

func (f *Factory[T]) onPaymentAlter(ctx context.Context, payload *payment.AlterPayload) error {
	if payload == nil {
		hlog.TraceFix("hyper.order.onPaymentAlter", ctx, fmt.Errorf("payload is nil"))
		return nil
	}
	switch payload.Dst {
	case payment.Initialized, payment.Executing:
		return nil
	default:
	}
	ordID := cast.ToUint64(payload.BizID)
	if ordID == 0 {
		hlog.TraceFix("hyper.order.onPaymentAlter", ctx,
			fmt.Errorf("payload.BizID is not a valid order id: %s", payload.BizID))
		return nil
	}
	eng, err := f.Load(ctx, ordID)
	if err != nil {
		return err
	}
	if eng == nil {
		return nil
	}
	switch payload.Dst {
	case payment.Completed:
		err = eng.doAdvToConsensus(ctx, ordID)
	case payment.Timeout:
		err = eng.doAdvToTimeout(ctx, ordID)
	case payment.Canceled:
		err = eng.doAdvToCanceled(ctx, ordID)
	default:
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}
