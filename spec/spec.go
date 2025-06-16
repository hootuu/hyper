package spec

import (
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/storage/hpg"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/category"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func ByCategory(cat category.ID) ([]*Spec, error) {
	arrM, err := hpg.Find[CatSpecM](func() *gorm.DB {
		return zplt.HelixPgDB().PG().Where("category = ?", cat)
	})
	if err != nil {
		hlog.Err("hyper.spec.ByCategory", zap.Error(err))
		return nil, err
	}
	if len(arrM) == 0 {
		return []*Spec{}, nil
	}
	var arr []*Spec
	for _, item := range arrM {
		arr = append(arr, item.To())
	}
	return arr, nil
}
