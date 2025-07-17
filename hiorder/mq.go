package hiorder

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/unicom/hmq/hmq"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hsys"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyper/hyperplt"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

var gFactoryMap = make(map[Code]func(ordID ID, payload *PaymentPayload) error)

func doRegFactoryMqHandle(code Code, handle func(ordID ID, payload *PaymentPayload) error) {
	gFactoryMap[code] = handle
}

func onPaymentAltered(ctx context.Context, msg *hmq.Message) (err error) {
	if hlog.IsElapseFunction() {
		defer hlog.ElapseWithCtx(ctx, "onPaymentAltered",
			hlog.F(zap.String("msg.id", msg.ID)),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err)}
				}
				return nil
			},
		)()
	}
	if msg == nil {
		hlog.Fix("hiorder.onPaymentAlter: msg is nil")
		return nil
	}
	payload, err := hjson.FromBytes[PaymentPayload](msg.Payload)
	if err != nil {
		hlog.Fix("hiorder.onPaymentAlter: payload is invalid")
		return nil
	}
	if payload == nil {
		hlog.Fix("hiorder.onPaymentAlter: payload is nil")
		return nil
	}
	orderCollar, err := collar.FromID(payload.OrderCollar)
	if err != nil {
		hlog.Fix("hiorder.onPaymentAlter: collar is invalid",
			zap.String("payload.OrderCollar", payload.OrderCollar),
			zap.Error(err))
		return nil
	}
	ordCode, ordIdStr := orderCollar.Parse()
	handle, ok := gFactoryMap[ordCode]
	if !ok {
		hlog.Fix("hiorder.onPaymentAlter: order code not registered", zap.String("code", ordCode))
		return nil
	}
	ordID := cast.ToUint64(ordIdStr)
	err = handle(ordID, payload)
	if err != nil {
		return err
	}
	return nil
}

var gOrderConsumer *hmq.Consumer

func init() {
	helix.AfterStartup(func() {
		fmt.Println("onPaymentAltered init.......") //todo
		gOrderConsumer = hyperplt.MQ().NewConsumer(
			"HYPER_ORD_PAY_ALTER_LISTENER",
			PaymentAlterTopic,
			"HIORDER",
		).WithHandler(onPaymentAltered)
		err := hyperplt.MQ().RegisterConsumer(gOrderConsumer)
		if err != nil {
			hsys.Exit(err)
		}
		fmt.Println("onPaymentAltered init.......[ok]") //todo
	})
}
