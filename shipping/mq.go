package shipping

import (
	"context"
	"github.com/hootuu/helix/unicom/hmq/hmq"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hretry"
	"github.com/hootuu/hyper/hyperplt"
	"go.uber.org/zap"
)

const (
	MqTopicShippingAlter = "HYPER_SHIPPING_ALTER"
)

type AlterPayload struct {
	ShippingID ID     `json:"shipping_id"`
	BizCode    string `json:"biz_code"`
	BizID      string `json:"biz_id"`
	Src        Status `json:"src"`
	Dst        Status `json:"dst"`
}

type AlterHandle func(ctx context.Context, payload *AlterPayload) error

func mqPublishShippingAlter(payload *AlterPayload) {
	var err error
	err = hretry.Must(func() error {
		return hyperplt.MqPublish(MqTopicShippingAlter, hjson.MustToBytes(payload))
	})
	if err != nil {
		hlog.Fix("hyper.shipping.notify", zap.Error(err),
			zap.Uint64("sid", payload.ShippingID),
			zap.String("biz_code", payload.BizCode),
			zap.String("biz_id", payload.BizID),
			zap.Int("src", int(payload.Src)),
			zap.Int("dst", int(payload.Dst)))
		return
	}
}

func onMqShippingAlter(ctx context.Context, msg *hmq.Message) error {
	if msg == nil {
		hlog.Fix("hyper.shipping.notify: msg is nil")
		return nil
	}
	payload := hjson.MustFromBytes[AlterPayload](msg.Payload)
	if payload == nil {
		hlog.Fix("hyper.shipping.notify: payload is nil")
		return nil
	}
	err := onShippingAlter(ctx, payload)
	if err != nil {
		return err
	}
	return nil
}

func doMqInit() error {
	mqShippingConsumer := hyperplt.MQ().NewConsumer(
		"hyper_shipping_mq",
		MqTopicShippingAlter,
		"HYPER",
	).WithHandler(onMqShippingAlter)
	err := hyperplt.MQ().RegisterConsumer(mqShippingConsumer)
	if err != nil {
		return err
	}
	return nil
}
