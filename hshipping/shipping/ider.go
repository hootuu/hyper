package shipping

import (
	"github.com/hootuu/helix/components/hnid"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hsys"
	"go.uber.org/zap"
	"sync"
)

var gSyncMu sync.Mutex
var gShippingIdGenerator hnid.Generator

func NxtShippingID() ID {
	if gShippingIdGenerator != nil {
		return gShippingIdGenerator.NextUint64()
	}
	gSyncMu.Lock()
	defer gSyncMu.Unlock()
	if gShippingIdGenerator != nil {
		return gShippingIdGenerator.NextUint64()
	}
	var err error
	gShippingIdGenerator, err = hnid.NewGenerator("hyper_shipping_id",
		hnid.NewOptions(1, 8).
			SetTimestamp(hnid.Minute, false).
			SetAutoInc(4, 1, 9999, 500),
	)
	if err != nil {
		hlog.Fix("hpay.ider.ShippingIdGenerator", zap.Error(err))
		hsys.Exit(err)
	}
	return gShippingIdGenerator.NextUint64()
}
