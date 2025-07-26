package shipping

import (
	"context"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/unicom/hmq/hmq"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hretry"
	"github.com/hootuu/hyle/hsys"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyper/hyperplt"
	"go.uber.org/zap"
)

const (
	MqTopicShipping = "HYPER_SHIPPING_ALTER"
)

type AlterPayload struct {
	ShippingID ID     `json:"shipping_id"`
	BizCode    string `json:"biz_code"`
	BizID      string `json:"biz_id"`
	Src        Status `json:"src"`
	Dst        Status `json:"dst"`
}

func (p *AlterPayload) IsCompleted() bool {
	return p.Dst == StatusDelivered
}

func mqPublishShippingAlter(sid ID, bizLink collar.Link, srcStatus Status, dstStatus Status) {
	var err error
	payload := &AlterPayload{
		ShippingID: sid,
		Src:        srcStatus,
		Dst:        dstStatus,
	}
	payload.BizCode, payload.BizID, err = bizLink.ToCodeID()
	if err != nil {
		hlog.Fix("hyper.shipping.notify: bizLink invalid", zap.Error(err))
		return
	}
	err = hretry.Must(func() error {
		return hyperplt.MqPublish(MqTopicShipping, hjson.MustToBytes(payload))
	})
	if err != nil {
		hlog.Fix("hyper.shipping.notify", zap.Error(err))
		return
	}
}

var gMqShippingConsumer *hmq.Consumer
var gMqShippingHandlerMap = make(map[string]func(ctx context.Context, payload *AlterPayload) error)

func MqRegisterShippingAlter(
	bizCode string,
	handle func(ctx context.Context, payload *AlterPayload) error,
) {
	gMqShippingHandlerMap[bizCode] = handle
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
	handle, ok := gMqShippingHandlerMap[payload.BizCode]
	if !ok {
		hlog.Fix("hyper.shipping.notify: handler not found", zap.String("biz_code", payload.BizCode))
		return nil
	}

	err := handle(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	helix.AfterStartup(func() {
		gMqShippingConsumer = hyperplt.MQ().NewConsumer(
			"hyper.shipping.mq",
			MqTopicShipping,
			"HYPER",
		).WithHandler(onMqShippingAlter)
		err := hyperplt.MQ().RegisterConsumer(gMqShippingConsumer)
		if err != nil {
			hsys.Exit(err)
		}
	})
}
