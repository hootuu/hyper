package pwh

import "github.com/hootuu/helix/components/hnid"

var gPwhIdGenerator hnid.Generator

func initPwhIdGenerator() error {
	var err error
	gPwhIdGenerator, err = hnid.NewGenerator("hyper_prod_pwh_id",
		hnid.NewOptions(1, 8).
			SetTimestamp(hnid.Hour, false).
			SetAutoInc(6, 1, 999999, 1000),
	)
	if err != nil {
		return err
	}
	return nil
}

func nextID() ID {
	return gPwhIdGenerator.NextUint64()
}
