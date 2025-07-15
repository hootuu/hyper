package brand

import (
	"github.com/hootuu/helix/components/hnid"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
)

var gIdGenerator hnid.Generator

func NextID() ID {
	return gIdGenerator.NextUint64()
}

func doInitIdGenerator() error {
	var err error
	gIdGenerator, err = hnid.NewGenerator("hyper_brand_ider",
		hnid.NewOptions(1, 8).
			SetTimestamp(hnid.Hour, false).
			SetAutoInc(5, 1, 99999, 1000),
	)
	if err != nil {
		hlog.Err("hyper.brand.init", zap.Error(err))
		return err
	}
	return nil
}
