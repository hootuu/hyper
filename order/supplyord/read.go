package supplyord

import (
	"context"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyper/hiorder"
	"github.com/hootuu/hyper/hyperplt"
)

func GetByProdOrderID(ctx context.Context, prodOrderID hiorder.ID) (*hiorder.Order[Matter], error) {
	return hdb.Get[hiorder.Order[Matter]](hyperplt.DB(), "link = ?", GetFactory().core.OrderCollar(prodOrderID).Link())
}
