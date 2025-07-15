package hyperidx

import (
	"context"
	"errors"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/storage/hcanal"
	"github.com/hootuu/helix/storage/hmeili"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/hiorder"
	"github.com/hootuu/hyper/hiprod/vwh"
	"github.com/hootuu/hyper/hyperidx/prodidx"
	"github.com/hootuu/hyper/hyperplt"
	"go.uber.org/zap"
	"gorm.io/gorm/schema"
)

func init() {
	helix.Use(helix.BuildHelix("hyper_index", func() (context.Context, error) {
		err := doInit()
		if err != nil {
			return nil, err
		}
		return nil, nil
	}, func(ctx context.Context) {

	}))
}

func doInit() error {
	meiliPtr := hyperplt.Meili()
	mArr := []schema.Tabler{
		&vwh.VirtualWhSkuM{},
		&hiorder.OrderM{},
	}
	idxArr := []hmeili.Indexer{
		&prodidx.TxVwhProdIndexer{},
		&TxOrdIndexer{},
	}
	for i, m := range mArr {
		if len(idxArr) <= i {
			hlog.Err("hyper.idx.init: idxArr.len < mArr.len")
			return errors.New("hyper.idx.init: idxArr.len < mArr.len")
		}
		err := hmeili.InitIndexer(meiliPtr, idxArr[i])
		if err != nil {
			hlog.Err("hyper.idx: init indexer failed",
				zap.String("idx", idxArr[i].GetName()), zap.Error(err))
			return err
		}
		hyperplt.Canal().RegisterAlterHandler(
			hcanal.NewIndexHandler(
				m.TableName(),
				idxArr[i],
				meiliPtr,
			),
		)
	}
	return nil
}
