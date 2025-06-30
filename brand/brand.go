package brand

import (
	"errors"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/idx"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
)

func addBrand(brM *BrM) (*BrM, error) {
	if brM.Name == "" {
		return nil, errors.New("require Name")
	}
	brM.ID = idx.New()
	err := hdb.Create[BrM](zplt.HelixPgDB().PG(), brM)
	if err != nil {
		hlog.Err("hyper.brand.addBrand: Create", zap.Error(err))
		return nil, err
	}
	return brM, nil
}
