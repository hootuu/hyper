package prod

import "github.com/hootuu/helix/components/hnid"

var gSpuIdGenerator hnid.Generator

func initSpuIdGenerator() error {
	var err error
	gSpuIdGenerator, err = hnid.NewGenerator("hyper_prod_spu_ider",
		hnid.NewOptions(1, 8).
			SetTimestamp(hnid.Minute, false).
			SetAutoInc(5, 1, 99999, 500),
	)
	if err != nil {
		return err
	}
	return nil
}

func nextSpuID() SpuID {
	return gSpuIdGenerator.NextUint64()
}

var gSkuIdGenerator hnid.Generator

func initSkuIdGenerator() error {
	var err error
	gSkuIdGenerator, err = hnid.NewGenerator("hyper_prod_sku_ider",
		hnid.NewOptions(1, 8).
			SetTimestamp(hnid.Minute, false).
			SetAutoInc(6, 1, 999999, 500),
	)
	if err != nil {
		return err
	}
	return nil
}

func nextSkuID() SkuID {
	return gSkuIdGenerator.NextUint64()
}

var gSpecOptionIdGenerator hnid.Generator

func initSpecOptIdGenerator() error {
	var err error
	gSpecOptionIdGenerator, err = hnid.NewGenerator("hyper_prod_spec_opt_id",
		hnid.NewOptions(1, 3).
			SetTimestamp(hnid.Minute, false).
			SetAutoInc(6, 1, 999999, 1000),
	)
	if err != nil {
		return err
	}
	return nil
}

func nextSpecOptID() SpecOptID {
	return gSpecOptionIdGenerator.NextUint64()
}
