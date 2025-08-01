package hiorder

import (
	"context"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyper/hyperplt"
)

func DbMustGet(ctx context.Context, orderId string) (*OrderM, error) {
	return hdb.MustGet[OrderM](hyperplt.Tx(ctx), "id = ?", orderId)
}
