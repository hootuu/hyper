package shipping

import (
	"github.com/hootuu/helix/components/hnid"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hsys"
	"go.uber.org/zap"
)

var gShippingIdGenerator hnid.Generator

func nxtShippingID() ID {
	if gShippingIdGenerator == nil {
		helix.OnceLoad("hyper_shipping_ider", func() {
			var err error
			gShippingIdGenerator, err = hnid.NewGenerator("hyper_shipping_ider",
				hnid.NewOptions(1, 8).
					SetTimestamp(hnid.Minute, false).
					SetAutoInc(4, 1, 9999, 500),
			)
			if err != nil {
				hlog.Fix("hyper.shipping.NxtShippingID", zap.Error(err))
				hsys.Exit(err)
			}
		})
	}
	return gShippingIdGenerator.NextUint64()
}
