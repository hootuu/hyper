package shipping

import (
	"context"
	"errors"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyle/hypes/ex"
	"github.com/hootuu/hyper/hyperplt"
	"go.uber.org/zap"
)

func AdvPrepared(
	ctx context.Context,
	uniLink collar.Link,
	courierCode CourierCode,
	trackingNo string,
	_ ex.Meta,
) (err error) {
	InitIfNeeded()
	if uniLink == "" {
		return errors.New("uni_link is required")
	}
	if courierCode == "" {
		return errors.New("courier_code is required")
	}
	if trackingNo == "" {
		return errors.New("tracking_no is required")
	}
	tx := hyperplt.Tx(ctx)
	shipM, err := hdb.Get[ShipM](tx, "uni_link = ?", uniLink)
	if err != nil {
		return err
	}
	if shipM == nil {
		return errors.New("uni_link not found: " + uniLink.Display())
	}
	return doAdvToSubmitted(ctx, shipM.ID, courierCode, trackingNo)
}

func AdvCompleted(
	ctx context.Context,
	uniLink collar.Link,
	_ ex.Meta,
) (err error) {
	InitIfNeeded()
	if uniLink == "" {
		return errors.New("uni_link is required")
	}
	if hlog.IsElapseComponent() {
		defer hlog.ElapseWithCtx(ctx, "hyper.shipping.Completed",
			hlog.F(zap.String("uni_link", uniLink.Str())),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{
						zap.String("uni_link", uniLink.Display()),
						zap.Error(err)}
				}
				return nil
			},
		)()
	}
	tx := hyperplt.Tx(ctx)
	shipM, err := hdb.Get[ShipM](tx, "uni_link = ?", uniLink)
	if err != nil {
		return err
	}
	if shipM == nil {
		return errors.New("uni_link not found: " + uniLink.Display())
	}
	return doAdvToCompleted(ctx, shipM.ID)
}
