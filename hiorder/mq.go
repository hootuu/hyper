package hiorder

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/unicom/hmq/hmq"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hretry"
	"github.com/hootuu/hyper/hyperplt"
	"go.uber.org/zap"
)

const (
	MqTopicOrderAlter = "HYPER_ORDER_ALTER"
)

type AlterPayload struct {
	OrderID ID     `json:"order_id"`
	Code    string `json:"code"`
	Src     Status `json:"src"`
	Dst     Status `json:"dst"`
}

type AlterHandle func(ctx context.Context, payload *AlterPayload) error

func mqPublishOrderAlter(payload *AlterPayload) {
	var err error
	err = hretry.Must(func() error {
		return hyperplt.MqPublish(MqTopicOrderAlter, hjson.MustToBytes(payload))
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

func onMqOrderAlter(ctx context.Context, msg *hmq.Message) error {
	if msg == nil {
		hlog.TraceFix("hyper.order.notify", ctx, fmt.Errorf("msg is nil"))
		return nil
	}
	fmt.Println("onMqOrderAlter: ", string(msg.Payload))

	payload := hjson.MustFromBytes[AlterPayload](msg.Payload)
	if payload == nil {
		hlog.TraceFix("hyper.order.notify", ctx, fmt.Errorf("payload is nil"))
		return nil
	}
	err := onUniOrdAlter(ctx, payload)
	if err != nil {
		return err
	}
	return nil
}

func doMqInit() error {
	mqShippingConsumer := hyperplt.MQ().NewConsumer(
		"hyper_order_mq",
		MqTopicOrderAlter,
		"HYPER",
	).WithHandler(onMqOrderAlter)
	err := hyperplt.MQ().RegisterConsumer(mqShippingConsumer)
	if err != nil {
		return err
	}
	return nil
}
