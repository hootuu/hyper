package pwh

import "github.com/hootuu/helix/components/hnid"

var gPwhIdGenerator hnid.Generator

func initPwhIdGenerator() error {
	var err error
	gPwhIdGenerator, err = hnid.NewGenerator("hyper_pwh_id",
		hnid.NewOptions(1, 6).
			SetTimestamp(hnid.Minute, false).
			SetAutoInc(6, 1, 999999, 1000),
	)
	if err != nil {
		return err
	}
	return nil
}
