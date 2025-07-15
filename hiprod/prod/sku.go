package prod

import (
	"context"
	"errors"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyper/hyperplt"
	"gorm.io/gorm"
)

func CreateSku(ctx context.Context, setting *SkuSpecSetting) (SkuID, error) {
	if setting.Spu == 0 {
		return 0, errors.New("require Spu")
	}
	skuID := nextSkuID()
	err := hdb.Tx(hyperplt.Tx(ctx), func(tx *gorm.DB) error {
		err := hdb.Create[SkuM](tx, &SkuM{
			ID:  skuID,
			Spu: setting.Spu,
		})
		if err != nil {
			return err
		}
		if len(setting.Specs) == 0 {
			return nil
		}
		for _, item := range setting.Specs {
			err := hdb.Create[SkuSpecM](tx, &SkuSpecM{
				Sku:     skuID,
				SpecOpt: item,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return skuID, nil
}
