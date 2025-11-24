package hiorder

import (
	"context"
	"errors"
	"fmt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/hyperplt"
	"github.com/hootuu/hyper/shipping"
	"github.com/spf13/cast"
	"time"
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
		err = eng.doAdvToCanceled(ctx, nil)
	default:
		return nil
	}
	return nil
}

func UpdateShipping(ctx context.Context, orderId, courierCode, trackingNo string) error {
	if orderId == "" {
		return errors.New("order_id is required")
	}
	if courierCode == "" {
		return errors.New("courier_code is required")
	}
	if trackingNo == "" {
		return errors.New("tracking_no is required")
	}
	err := hdb.Update[shipping.ShipM](hyperplt.Tx(ctx), map[string]any{
		"courier_code": courierCode,
		"tracking_no":  trackingNo,
	}, "biz_id = ?", orderId)
	if err != nil {
		return err
	}
	return hdb.Update[OrderM](hyperplt.Tx(ctx), map[string]any{
		"updated_at": time.Now(),
	}, "id = ?", orderId)
}

func UpdateShippingAddr(ctx context.Context, params shipping.UpdateAddrParams) error {
	if params.OrderId == "" {
		return errors.New("order_id is required")
	}
	orderM, err := DbMustGet(ctx, params.OrderId)
	if err != nil {
		return err
	}
	if orderM.Status != Consensus {
		return errors.New("order status is not consensus, can not update shipping address")
	}
	return shipping.UpdateAddrInfo(ctx, params)
}
