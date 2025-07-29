package shipping

import (
	"context"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyper/hyperplt"
)

func callbackShipCheckTimeoutAndAdv(ctx context.Context, shippingID ID) error {
	shipM, err := hdb.Get[ShipM](hyperplt.Tx(ctx), "id = ?", shippingID)
	if err != nil {
		return err
	}
	if shipM == nil {
		return nil
	}
	switch shipM.Status {
	case Completed:
		return nil
	case Submitted:
		return doAdvToTimeoutCompleted(ctx, shippingID)
	default:
		return nil
	}
}
