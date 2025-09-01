package supplyord

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/ticktock"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hretry"
	"github.com/hootuu/hyper/hiorder"
	"github.com/hootuu/hyper/hyperplt"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"time"
)

const (
	ttOrderTimeout = "PROD_ORD_ORDER_TIMEOUT"
)

func ttListenOrderTimeout(ctx context.Context, orderID hiorder.ID, timeout time.Duration) {
	if timeout == 0 {
		return
	}
	if orderID == 0 {
		hlog.TraceFix("prodord.order.ttListenOrderTimeout: orderID is 0", ctx, nil)
		return
	}
	hretry.Universal(func() error {
		err := hyperplt.Postman().Send(ctx, &ticktock.DelayJob{
			Type:      ttOrderTimeout,
			ID:        fmt.Sprintf("PROD_ORD:ORDER:TIMEOUT:%d", orderID),
			Payload:   []byte(fmt.Sprintf("%d", orderID)),
			UniqueTTL: 0,
			Delay:     timeout,
		})
		if err != nil {
			hlog.TraceFix("send ticktock:PROD_ORD_ORDER_TIMEOUT failed", ctx, err, zap.Uint64("orderID", orderID))
			return err
		}
		return err
	})
}

func init() {
	helix.Ready(func() {
		hyperplt.Ticktock().RegisterJobHandlerFunc(
			ttOrderTimeout,
			func(ctx context.Context, job *ticktock.Job) error {
				if job.Payload == nil {
					return nil
				}
				orderIdStr := string(job.Payload)
				if orderIdStr == "" {
					return nil
				}
				orderId := cast.ToUint64(orderIdStr)
				if orderId == 0 {
					return nil
				}
				return callbackOrderCheckTimeoutAndAdv(ctx, orderId)
			},
		)
	})
}
