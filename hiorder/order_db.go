package hiorder

import (
	"context"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyper/hyperplt"
)

func DbMustGet(ctx context.Context, orderId string) (*OrderM, error) {
	return hdb.MustGet[OrderM](hyperplt.Tx(ctx), "id = ?", orderId)
}

func DbExists(ctx context.Context, link collar.Link) (bool, error) {
	var count int64
	err := hyperplt.Tx(ctx).Model(&OrderM{}).Where("link = ?", link).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
