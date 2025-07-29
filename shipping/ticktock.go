package shipping

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/ticktock"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hretry"
	"github.com/hootuu/hyper/hyperplt"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"time"
)

const (
	ttShippingTimeout = "HYPER_SHIPPING_TIMEOUT"
)

func ttListenShippingTimeout(ctx context.Context, shippingID ID, timeout time.Duration) {
	if timeout == 0 {
		return
	}
	hretry.Universal(func() error {
		err := hyperplt.Postman().Send(ctx, &ticktock.DelayJob{
			Type:      ttShippingTimeout,
			ID:        fmt.Sprintf("hyper_shipping:%d", shippingID),
			Payload:   []byte(fmt.Sprintf("%d", shippingID)),
			UniqueTTL: 0,
			Delay:     timeout,
		})
		if err != nil {
			hlog.TraceFix("send ticktock:HYPER_SHIPPING_TIMEOUT failed", ctx, err, zap.Uint64("shippingID", shippingID))
			return err
		}
		return err
	})
}

func init() {
	helix.Ready(func() {
		hyperplt.Ticktock().RegisterJobHandlerFunc(
			ttShippingTimeout,
			func(ctx context.Context, job *ticktock.Job) error {
				if job.Payload == nil {
					return nil
				}
				shippingIdStr := string(job.Payload)
				if shippingIdStr == "" {
					return nil
				}
				shippingId := cast.ToUint64(shippingIdStr)
				if shippingId == 0 {
					return nil
				}
				return callbackShipCheckTimeoutAndAdv(ctx, shippingId)
			},
		)
	})
}
