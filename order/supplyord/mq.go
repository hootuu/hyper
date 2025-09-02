package supplyord

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/unicom/hmq/hmq"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/hiorder"
	"github.com/hootuu/hyper/hyperplt"
	ordtypes "github.com/hootuu/hyper/order/types"
	"go.uber.org/zap"
)

func init() {
	helix.AfterStartup(func() {
		err := doMqInit()
		if err != nil {
			hlog.Fix("hyper.supplyord.notify", zap.Error(err))
		}
	})
}

func doMqInit() error {
	mqOrderConsumer := hyperplt.MQ().NewConsumer(
		"hyper.supplyord.mq",
		ordtypes.MqTopicProdOrderAlter,
		"SUPPLYORD",
	).WithHandler(onMqOrderAlter)
	err := hyperplt.MQ().RegisterConsumer(mqOrderConsumer)
	if err != nil {
		return err
	}
	return nil
}

func onMqOrderAlter(ctx context.Context, msg *hmq.Message) error {
	if msg == nil {
		hlog.TraceFix("hyper.supplyord.notify", ctx, fmt.Errorf("msg is nil"))
		return nil
	}

	payload := hjson.MustFromBytes[ordtypes.AlterPayload](msg.Payload)
	if payload == nil {
		hlog.TraceFix("hyper.supplyord.notify", ctx, fmt.Errorf("payload is nil"))
		return nil
	}

	switch payload.Dst {
	case hiorder.Consensus:
		return doOrderConsensusAdv(ctx, payload.OrderID)
	case hiorder.Executing:
		return doOrderExecutingAdv(ctx, payload.OrderID)
	case hiorder.Completed:
		return doOrderCompletedAdv(ctx, payload.OrderID)
	default:
		return nil
	}
}
