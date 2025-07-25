package prod

import (
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyper/hyperplt"
)

func DbSpuGet(id SpuID) (*SpuM, error) {
	return hdb.Get[SpuM](hyperplt.DB(), "id = ?", id)
}
