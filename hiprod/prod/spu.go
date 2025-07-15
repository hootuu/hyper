package prod

import (
	"context"
	"errors"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/hyperplt"
	"go.uber.org/zap"
)

func CreateSpu(ctx context.Context, spuM *SpuM) (*SpuM, error) {
	if spuM.Name == "" {
		return nil, errors.New("require Name")
	}
	spuM.ID = nextSpuID()
	err := hdb.Create[SpuM](hyperplt.Tx(ctx), spuM)
	if err != nil {
		hlog.Err("hyper.prod.CreateSpu", zap.Error(err))
		return nil, err
	}
	return spuM, nil
}
