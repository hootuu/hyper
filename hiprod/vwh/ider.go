package vwh

import "github.com/hootuu/helix/components/hnid"

var gVwhIdGenerator hnid.Generator

func initVwhIdGenerator() error {
	var err error
	gVwhIdGenerator, err = hnid.NewGenerator("hyper_vwh_id",
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
	return gVwhIdGenerator.NextUint64()
}
