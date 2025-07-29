package shipping

import (
	"context"
	"errors"
	"fmt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/ex"
	"github.com/hootuu/hyper/hyperplt"
	"go.uber.org/zap"
)

func AdvPrepared(
	ctx context.Context,
	id ID,
	courierCode CourierCode,
	trackingNo string,
	_ ex.Meta,
) (err error) {
	InitIfNeeded()
	if id == 0 {
		return errors.New("id is required")
	}
	if courierCode == "" {
		return errors.New("courier_code is required")
	}
	if trackingNo == "" {
		return errors.New("tracking_no is required")
	}
	tx := hyperplt.Tx(ctx)
	shipM, err := hdb.Get[ShipM](tx, "id = ?", id)
	if err != nil {
		return err
	}
	if shipM == nil {
		return fmt.Errorf("id not found: %d", id)
	}
	return doAdvToSubmitted(ctx, shipM.ID, courierCode, trackingNo)
}

func AdvCompleted(
	ctx context.Context,
	id ID,
	_ ex.Meta,
) (err error) {
	InitIfNeeded()
	if id == 0 {
		return errors.New("id is required")
	}
	if hlog.IsElapseComponent() {
		defer hlog.ElapseWithCtx(ctx, "hyper.shipping.Completed",
			hlog.F(zap.Uint64("id", id)),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{
						zap.Uint64("id", id),
						zap.Error(err)}
				}
				return nil
			},
		)()
	}
	tx := hyperplt.Tx(ctx)
	shipM, err := hdb.Get[ShipM](tx, "id = ?", id)
	if err != nil {
		return err
	}
	if shipM == nil {
		return fmt.Errorf("uni_link not found: %d", id)
	}
	return doAdvToCompleted(ctx, shipM.ID)
}
