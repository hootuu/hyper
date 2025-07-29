package payment

import (
	"github.com/hootuu/helix/components/hnid"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hsys"
	"go.uber.org/zap"
	"sync"
)

var gSyncMu sync.Mutex
var gPaymentIdGenerator hnid.Generator

func NxtPaymentID() ID {
	gPaymentIdGenerator = func() hnid.Generator {
		gSyncMu.Lock()
		defer gSyncMu.Unlock()
		if gPaymentIdGenerator != nil {
			return gPaymentIdGenerator
		}
		var err error
		gPaymentIdGenerator, err = hnid.NewGenerator("hyper_payment_id",
			hnid.NewOptions(1, 8).
				SetTimestamp(hnid.Second, false).
				SetAutoInc(5, 1, 99999, 1000),
		)
		if err != nil {
			hlog.Fix("hpay.ider.PaymentIdGenerator", zap.Error(err))
			hsys.Exit(err)
		}
		return gPaymentIdGenerator
	}()
	return gPaymentIdGenerator.NextUint64()
}
