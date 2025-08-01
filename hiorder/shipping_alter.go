package hiorder

import (
	"context"
	"fmt"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/shipping"
	"github.com/spf13/cast"
)

func (f *Factory[T]) onShippingAlter(ctx context.Context, payload *shipping.AlterPayload) error {
	if payload == nil {
		hlog.TraceFix("hyper.order.onShippingAlter", ctx, fmt.Errorf("payload is nil"))
		return nil
	}
	switch payload.Dst {
	case shipping.Initialized:
		return nil
	default:
	}
	ordID := cast.ToUint64(payload.BizID)
	if ordID == 0 {
		hlog.TraceFix("hyper.order.onShippingAlter", ctx,
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
	case shipping.Submitted:
		err = eng.doAdvToExecuting(ctx)
	case shipping.Completed:
		err = eng.doAdvToCompleted(ctx)
	case shipping.Canceled:
		err = eng.doAdvToCanceled(ctx)
	default:
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}
