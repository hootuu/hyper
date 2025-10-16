package supplyord

import (
	"context"
	"errors"
	"fmt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/hiorder"
	"github.com/hootuu/hyper/hyperplt"
	"github.com/spf13/cast"
	"time"
)

func DoOrderRefund(ctx context.Context, orderID hiorder.ID) error {
	orderM, err := hiorder.DbMustGet(ctx, cast.ToString(orderID))
	if err != nil {
		return nil
	}
	if orderM.Status != hiorder.Executing && orderM.Status != hiorder.Consensus {
		return errors.New("order status is invalid")
	}
	err = hdb.Update[hiorder.OrderM](hyperplt.DB(), map[string]any{
		"status":         hiorder.Refunded,
		"completed_time": time.Now(),
	}, "id = ?", orderID)
	if err != nil {
		return err
	}
	return nil
}

func DoOrderRefundAdv(ctx context.Context, orderID hiorder.ID) error {
	orderM, err := hiorder.DbMustGet(ctx, cast.ToString(orderID))
	if err != nil {
		return err
	}
	if orderM.Status == hiorder.Refunded {
		supOrder, err := GetByProdOrderID(ctx, orderID)
		if err != nil {
			return err
		}
		if supOrder == nil {
			hlog.TraceFix("hyper.supplyord.doOrderCompletedAdv", ctx, fmt.Errorf("supply order not found: %d", orderID))
			return nil
		}
		err = hdb.Update[hiorder.OrderM](hyperplt.DB(), map[string]any{
			"status":         hiorder.Refunded,
			"completed_time": orderM.CompletedTime,
		}, "id = ?", supOrder.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
