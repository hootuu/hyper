package supplyord

import (
	"context"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyper/hiorder"
	"github.com/hootuu/hyper/hyperplt"
	"github.com/spf13/cast"
)

func GetByProdOrderID(ctx context.Context, prodOrderID hiorder.ID) (*hiorder.OrderM, error) {
	return hdb.Get[hiorder.OrderM](hyperplt.DB(), "link = ?", GetFactory().core.OrderCollar(prodOrderID).Link())
}

func GetProdOrdID(ctx context.Context, supplyOrdId hiorder.ID) (hiorder.ID, error) {
	supOrd, err := hdb.MustGet[hiorder.OrderM](hyperplt.DB(), "id = ?", supplyOrdId)
	if err != nil {
		return 0, nil
	}
	meta := hjson.MustFromBytes[dict.Dict](supOrd.Meta)
	return cast.ToUint64(meta.Get("orderId").String()), nil
}
