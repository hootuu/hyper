package prod

import (
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyper/hyperplt"
)

func DbSkuGet(id SkuID) (*SkuM, error) {
	return hdb.Get[SkuM](hyperplt.DB(), "id = ?", id)
}
