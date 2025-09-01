package prodord

import (
	"context"
	"github.com/hootuu/hyper/hiorder"
	"github.com/spf13/cast"
)

func callbackOrderCheckTimeoutAndAdv(ctx context.Context, orderID hiorder.ID) error {
	orderM, err := hiorder.DbMustGet(ctx, cast.ToString(orderID))
	if err != nil {
		return err
	}
	if orderM == nil {
		return nil
	}
	if orderM.Status == Initial {
		e, err := NewFactory(orderM.Code).Core().Load(ctx, orderID)
		if err != nil {
			return err
		}
		return e.DoAdvToTimeout(ctx)
	}
	return nil
}
