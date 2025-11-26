package shipping

import (
	"context"
	"errors"
	"fmt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/hyperplt"
	"go.uber.org/zap"
	"time"
)

var gShippingHandlerMap map[string]AlterHandle

func doRegisterAlterHandle(
	bizCode string,
	handle AlterHandle,
) {
	if gShippingHandlerMap == nil {
		gShippingHandlerMap = make(map[string]AlterHandle)
	}
	gShippingHandlerMap[bizCode] = handle
}

func onShippingAlter(ctx context.Context, payload *AlterPayload) error {
	if gShippingHandlerMap == nil {
		return nil
	}
	shipM, err := hdb.Get[ShipM](hyperplt.Tx(ctx), "id = ? AND status = ?", payload.ShippingID, payload.Dst)
	if err != nil {
		return err
	}
	if shipM == nil {
		hlog.TraceFix("hyper.shipping.onShippingAlter", ctx,
			fmt.Errorf("no such shipping: %d", payload.ShippingID),
			zap.Uint64("shippingID", payload.ShippingID),
			zap.String("bizCode", payload.BizCode),
			zap.String("bizID", payload.BizID),
			zap.Int("src", int(payload.Src)),
			zap.Int("dst", int(payload.Dst)))
		return nil
	}

	switch payload.Dst {
	case Submitted:
		timeout := 7 * 24 * time.Hour
		if shipM.Timeout > 0 {
			timeout = shipM.Timeout
		}
		ttListenShippingTimeout(ctx, shipM.ID, timeout)
	default:
	}

	handle, ok := gShippingHandlerMap[shipM.BizCode]
	if !ok {
		hlog.TraceFix("hyper.shipping.notify",
			ctx, errors.New("no handler for this biz code"),
			zap.String("biz_code", shipM.BizCode))
		return nil
	}

	err = handle(ctx, payload)
	if err != nil {
		return err
	}
	return nil
}

type UpdateAddrParams struct {
	OrderId string `json:"orderId"`
	Name    string `json:"name"`
	Mobi    string `json:"mobi"`
	Address string `json:"address"`
}

func UpdateAddrInfo(ctx context.Context, params UpdateAddrParams) error {
	if params.OrderId == "" {
		return errors.New("orderId is required")
	}
	if params.Name == "" {
		return errors.New("name is required")
	}
	if params.Mobi == "" {
		return errors.New("mobi is required")
	}
	if params.Address == "" {
		return errors.New("address is required")
	}

	addrMap := map[string]any{
		"province": "",
		"city":     "",
		"district": "",
		"address":  params.Address,
		"contact": map[string]any{
			"mobi": params.Mobi,
			"name": params.Name,
		},
	}
	return hdb.Update[ShipM](hyperplt.Tx(ctx), map[string]any{
		"address": hjson.MustToBytes(addrMap),
	}, "biz_id = ?", params.OrderId)
}
