package product

import (
	"errors"
	"github.com/hootuu/helix/components/hnid"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/storage/hpg"
	"gorm.io/gorm"
)

var gSkuIdGenerator hnid.Generator

func initSkuIdGenerator() error {
	var err error
	gSkuIdGenerator, err = hnid.NewGenerator("hyper_sku_id",
		hnid.NewOptions(1, 8).
			SetTimestamp(hnid.Minute, false).
			SetAutoInc(9, 1, 999999999, 1000),
	)
	if err != nil {
		return err
	}
	return nil
}

//func CreateSku(skuM *SkuM) (*SkuM, error) {
//	//if skuM.Spu == "" {
//	//	return nil, errors.New("require Spu")
//	//}
//	//skuM.ID = gSkuIdGenerator.NextString()
//	//err := hpg.Create[SkuM](zplt.HelixPgDB().PG(), skuM)
//	//if err != nil {
//	//	hlog.Err("hyper.product.CreateSku", zap.Error(err))
//	//	return nil, err
//	//}
//	//return skuM, nil
//}

func CreateSku(setting *SkuSpecSetting) (SkuID, error) {
	return createSku(zplt.HelixPgDB().PG(), setting)
}

func createSku(tx *gorm.DB, setting *SkuSpecSetting) (SkuID, error) {
	if setting.Spu == "" {
		return "", errors.New("require Spu")
	}
	skuID := gSkuIdGenerator.NextString()
	err := hpg.Create[SkuM](tx, &SkuM{
		ID:  skuID,
		Spu: setting.Spu,
	})
	if err != nil {
		return "", err
	}
	if len(setting.Specs) == 0 {
		return "", nil
	}
	for _, item := range setting.Specs {
		err := hpg.Create[SkuSpecM](tx, &SkuSpecM{
			Sku:     skuID,
			SpecOpt: item,
		})
		if err != nil {
			return "", err
		}
	}
	return skuID, nil
}
