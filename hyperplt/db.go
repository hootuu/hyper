package hyperplt

import (
	"context"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/components/zplt/zmeili"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/helix/storage/hmeili"
	"gorm.io/gorm"
)

func Meili() *hmeili.Meili {
	return zmeili.HelixMeili()
}

func DB() *gorm.DB {
	return zplt.HelixPgDB().PG()
}

func Tx(ctx context.Context) *gorm.DB {
	tx := hdb.CtxTx(ctx)
	if tx == nil {
		tx = DB()
	}
	return tx
}

func Ctx(ctx ...context.Context) context.Context {
	return hdb.TxCtx(DB(), ctx...)
}
