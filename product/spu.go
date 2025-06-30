package product

import (
	"errors"
	"github.com/hootuu/helix/components/hnid"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var gSpuIdGenerator hnid.Generator

func initSpuIdGenerator() error {
	var err error
	gSpuIdGenerator, err = hnid.NewGenerator("hyper_spu_id",
		hnid.NewOptions(1, 6).
			SetTimestamp(hnid.Minute, false).
			SetAutoInc(6, 1, 999999, 1000),
	)
	if err != nil {
		return err
	}
	return nil
}

func CreateSpu(spuM *SpuM) (*SpuM, error) {
	return createSpu(zplt.HelixPgDB().PG(), spuM)
}

func createSpu(tx *gorm.DB, spuM *SpuM) (*SpuM, error) {
	if spuM.Name == "" {
		return nil, errors.New("require Name")
	}
	spuM.ID = gSpuIdGenerator.NextString()
	err := hdb.Create[SpuM](tx, spuM)
	if err != nil {
		hlog.Err("hyper.product.createSpu", zap.Error(err))
		return nil, err
	}
	return spuM, nil
}
