package hiorder

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/hyperplt"
	"go.uber.org/zap"
)

func (f *Factory[T]) RegisterAlterHandle(handle AlterHandle) {
	f.alterHandlerArr = append(f.alterHandlerArr, handle)
}

func (f *Factory[T]) onOrdAlter(ctx context.Context, payload *AlterPayload) error {
	if len(f.alterHandlerArr) == 0 {
		return nil
	}
	ordM, err := hdb.Get[OrderM](hyperplt.Tx(ctx), "id = ? AND status = ?",
		payload.OrderID, payload.Dst)
	if err != nil {
		return err
	}
	if ordM == nil {
		hlog.TraceFix("hyper.order.onOrdAlter", ctx,
			fmt.Errorf("no such order: %d", payload.OrderID),
			zap.Uint64("orderID", payload.OrderID),
			zap.Int("src", int(payload.Src)),
			zap.Int("dst", int(payload.Dst)))
		return nil
	}

	for _, handle := range f.alterHandlerArr {
		err = handle(ctx, &AlterPayload{
			OrderID: ordM.ID,
			Src:     payload.Src,
			Dst:     payload.Dst,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

var gUniOrdAlterHandlerMap map[string]AlterHandle

func doInjectUniOrdAlterHandler(code string, handle AlterHandle) {
	if gUniOrdAlterHandlerMap == nil {
		gUniOrdAlterHandlerMap = make(map[string]AlterHandle)
	}
	gUniOrdAlterHandlerMap[code] = handle
}

func onUniOrdAlter(ctx context.Context, payload *AlterPayload) error {
	if len(gUniOrdAlterHandlerMap) == 0 {
		return nil
	}
	alterHandle, ok := gUniOrdAlterHandlerMap[payload.Code]
	if !ok {
		hlog.TraceFix("hyper.order.onUniOrdAlter", ctx,
			fmt.Errorf("no such order alter handle: %s", payload.Code))
		return nil
	}
	return alterHandle(ctx, payload)
}
