package payment

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
	MqTopicPayment = "HYPER_PAYMENT_ALTER"
)

type AlterPayload struct {
	PaymentID ID     `json:"payment_id"`
	BizCode   string `json:"biz_code"`
	BizID     string `json:"biz_id"`
	Src       Status `json:"src"`
	Dst       Status `json:"dst"`
}

func (p *AlterPayload) IsCompleted() bool {
	return p.Dst == Completed
}

func mqPublishPaymentAlter(pid ID, bizLink collar.Link, srcStatus Status, dstStatus Status) {
	var err error
	payload := &AlterPayload{
		PaymentID: pid,
		Src:       srcStatus,
		Dst:       dstStatus,
	}
	payload.BizCode, payload.BizID, err = bizLink.ToCodeID()
	if err != nil {
		hlog.Fix("hyper.payment.notify: bizLink invalid", zap.Error(err))
		return
	}
	err = hretry.Must(func() error {
		return hyperplt.MqPublish(MqTopicPayment, hjson.MustToBytes(payload))
	})
	if err != nil {
		hlog.Fix("hyper.payment.notify", zap.Error(err))
		return
	}
}

var gMqPaymentConsumer *hmq.Consumer
var gMqPaymentHandlerMap = make(map[string]func(ctx context.Context, payload *AlterPayload) error)

func MqRegisterPaymentAlter(
	bizCode string,
	handle func(ctx context.Context, payload *AlterPayload) error,
) {
	gMqPaymentHandlerMap[bizCode] = handle
}

func onMqPaymentAlter(ctx context.Context, msg *hmq.Message) error {
	if msg == nil {
		hlog.Fix("hyper.payment.notify: msg is nil")
		return nil
	}
	payload := hjson.MustFromBytes[AlterPayload](msg.Payload)
	if payload == nil {
		hlog.Fix("hyper.payment.notify: payload is nil")
		return nil
	}
	handle, ok := gMqPaymentHandlerMap[payload.BizCode]
	if !ok {
		hlog.Fix("hyper.payment.notify: handler not found", zap.String("biz_code", payload.BizCode))
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
		gMqPaymentConsumer = hyperplt.MQ().NewConsumer(
			"hyper.payment.mq",
			MqTopicPayment,
			"HYPER",
		).WithHandler(onMqPaymentAlter)
		err := hyperplt.MQ().RegisterConsumer(gMqPaymentConsumer)
		if err != nil {
			hsys.Exit(err)
		}
	})
}
