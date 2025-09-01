package hyperidxsync

import (
	"errors"
	"github.com/hootuu/helix/storage/hcanal"
	"github.com/hootuu/helix/storage/hmeili"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hsys"
	"github.com/hootuu/hyper/hiorder"
	"github.com/hootuu/hyper/hiprod/prod"
	"github.com/hootuu/hyper/hiprod/pwh"
	"github.com/hootuu/hyper/hiprod/vwh"
	"github.com/hootuu/hyper/hyperidx"
	"github.com/hootuu/hyper/hyperidx/prodidx"
	"github.com/hootuu/hyper/hyperplt"
	"go.uber.org/zap"
	"gorm.io/gorm/schema"
)

func Use() {
	err := doInit()
	if err != nil {
		hsys.Error("init hyperidx syncer failed: " + err.Error())
		hsys.Exit(err)
	}
}

func doInit() error {
	meiliPtr := hyperplt.Meili()
	mArr := []schema.Tabler{
		&vwh.VirtualWhSkuM{},
		&hiorder.OrderM{},
		&pwh.PhysicalSkuM{},
		&prod.SpuM{},
	}
	idxArr := []hmeili.Indexer{
		&prodidx.TxVwhProdIndexer{},
		&hyperidx.TxOrdIndexer{},
		&prodidx.TxPwhProdIndexer{},
		&prodidx.TxSpuIndexer{},
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
		canal().RegisterAlterHandler(
			hcanal.NewIndexHandler(
				m.TableName(),
				idxArr[i],
				meiliPtr,
			),
		)
	}
	return nil
}
