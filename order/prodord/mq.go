package prodord

import (
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hretry"
	"github.com/hootuu/hyper/hyperplt"
	"github.com/hootuu/hyper/order/types"
	"go.uber.org/zap"
)

func mqPublishOrderAlter(payload *ordtypes.AlterPayload) {
	var err error
	err = hretry.Must(func() error {
		return hyperplt.MqPublish(ordtypes.MqTopicProdOrderAlter, hjson.MustToBytes(payload))
	})
	if err != nil {
		hlog.Fix("hyper.order.notify", zap.Error(err),
			zap.String("biz_code", payload.Code),
			zap.Uint64("id", payload.OrderID),
			zap.Int("src", int(payload.Src)),
			zap.Int("dst", int(payload.Dst)))
		return
	}
	return
}
