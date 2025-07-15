package prod

import (
	"errors"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func CreateSpec(setting *SpuSpecSetting) (*SpuSpecSetting, error) {
	return createSpec(zplt.HelixPgDB().PG(), setting)
}

func createSpec(tx *gorm.DB, setting *SpuSpecSetting) (*SpuSpecSetting, error) {
	if len(setting.Specs) == 0 {
		return nil, errors.New("require Specs")
	}
	for i, spuSpecItem := range setting.Specs {
		if spuSpecItem.Seq == 0 {
			spuSpecItem.Seq = i + 1
		}
		err := doCreateSpecItem(tx, setting.Spu, spuSpecItem)
		if err != nil {
			return nil, err
		}
	}
	return setting, nil
}

func doCreateSpecItem(tx *gorm.DB, spuID SpuID, item *SpuSpec) error {
	if len(item.Options) == 0 {
		return errors.New("require SpuSpec.Options")
	}
	err := hdb.Create[SpuSpecM](tx, &SpuSpecM{
		Spu:  spuID,
		Spec: item.Spec,
		Seq:  item.Seq,
	})
	if err != nil {
		hlog.Err("hyper.product.doCreateSpecItem: hdb.Create[SpuSpecM]", zap.Error(err))
		return err
	}

	for i, optItem := range item.Options {
		optItem.ID = nextSpecOptID()
		if optItem.Seq == 0 {
			optItem.Seq = i + 1
		}
		err := hdb.Create[SpuSpecOptM](tx, &SpuSpecOptM{
			ID:    optItem.ID,
			Spu:   spuID,
			Spec:  item.Spec,
			Label: optItem.Label,
			Media: hjson.MustToBytes(optItem.Media),
			Seq:   optItem.Seq,
		})
		if err != nil {
			hlog.Err("hyper.product.doCreateSpecItem: hdb.Create[SpuSpecOptM]", zap.Error(err))
			return err
		}
	}
	return nil
}
