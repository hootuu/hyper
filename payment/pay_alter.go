package payment

import (
	"context"
	"fmt"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
)

var gPaymentHandlerMap map[string]AlterHandle

func doRegisterPaymentAlter(
	bizCode string,
	handle AlterHandle,
) {
	if gPaymentHandlerMap == nil {
		gPaymentHandlerMap = make(map[string]AlterHandle)
	}
	gPaymentHandlerMap[bizCode] = handle
}

func onPaymentAlter(ctx context.Context, payload *AlterPayload) error {
	if len(gPaymentHandlerMap) == 0 {
		return nil
	}
	handle, ok := gPaymentHandlerMap[payload.BizCode]
	if !ok {
		hlog.TraceFix("hyper.payment.notify",
			ctx, fmt.Errorf("no handler for bizcode %s", payload.BizCode),
			zap.Any("handles", gPaymentHandlerMap),
			zap.String("biz_code", payload.BizCode),
		)
		return nil
	}

	err := handle(ctx, payload)
	if err != nil {
		return err
	}
	return nil
}
