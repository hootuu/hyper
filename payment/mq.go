package payment

import (
	"context"
	"github.com/hootuu/helix/unicom/hmq/hmq"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hretry"
	"github.com/hootuu/hyper/hyperplt"
)

const (
	MqTopicPaymentAlter    = "HYPER_PAYMENT_ALTER"
	MqTopicPaymentJobAlter = "HYPER_PAYMENT_JOB_ALTER"
)

type AlterPayload struct {
	PaymentID ID     `json:"payment_id"`
	BizCode   string `json:"biz_code"`
	BizID     string `json:"biz_id"`
	Src       Status `json:"src"`
	Dst       Status `json:"dst"`
}

type AlterHandle func(ctx context.Context, payload *AlterPayload) error

type JobAlterPayload struct {
	JobID JobID     `json:"job_id"`
	Src   JobStatus `json:"src"`
	Dst   JobStatus `json:"dst"`
}

func mqPublishPaymentAlter(ctx context.Context, payload *AlterPayload) {
	err := hretry.Must(func() error {
		return hyperplt.MqPublish(MqTopicPaymentAlter, hjson.MustToBytes(payload))
	})
	if err != nil {
		hlog.TraceFix("hyper.payment.notify", ctx, err)
		return
	}
}

func mqPublishJobAlter(ctx context.Context, payload *JobAlterPayload) {
	err := hretry.Must(func() error {
		return hyperplt.MqPublish(MqTopicPaymentJobAlter, hjson.MustToBytes(payload))
	})
	if err != nil {
		hlog.TraceFix("hyper.payment.job.notify", ctx, err)
		return
	}
}

var gMqPaymentConsumer *hmq.Consumer
var gMqPaymentJobConsumer *hmq.Consumer

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

	return onPaymentAlter(ctx, payload)
}

func onMqJobAlter(ctx context.Context, msg *hmq.Message) error {
	if msg == nil {
		hlog.Fix("hyper.payment.job.notify: msg is nil")
		return nil
	}
	payload := hjson.MustFromBytes[JobAlterPayload](msg.Payload)
	if payload == nil {
		hlog.Fix("hyper.payment.job.notify: payload is nil")
		return nil
	}
	return onJobAlter(ctx, payload.JobID, payload.Src, payload.Dst)
}

func doMqInit() error {
	gMqPaymentConsumer = hyperplt.MQ().NewConsumer(
		"hyper.payment.mq",
		MqTopicPaymentAlter,
		"HYPER",
	).WithHandler(onMqPaymentAlter)
	err := hyperplt.MQ().RegisterConsumer(gMqPaymentConsumer)
	if err != nil {
		return err
	}

	gMqPaymentJobConsumer = hyperplt.MQ().NewConsumer(
		"hyper.payment.job.mq",
		MqTopicPaymentJobAlter,
		"HYPER",
	).WithHandler(onMqJobAlter)
	err = hyperplt.MQ().RegisterConsumer(gMqPaymentJobConsumer)
	if err != nil {
		return err
	}
	return nil
}
