package shipping

import (
	"context"
	"errors"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyle/hypes/ex"
	"github.com/hootuu/hyper/hyperplt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CreateParas struct {
	UniLink collar.Link `json:"uni_link"`
	Address *Address    `json:"address"`
	Ex      *ex.Ex      `json:"ex"`
}

func (paras *CreateParas) Validate() error {
	if paras.UniLink == "" {
		return errors.New("uni_link is required")
	}
	if paras.Address == nil {
		return errors.New("address is required")
	}
	if err := paras.Address.Validate(); err != nil {
		return err
	}
	return nil
}

func Create(ctx context.Context, paras *CreateParas) (id ID, err error) {
	if paras == nil {
		return 0, errors.New("paras is required")
	}
	if err := paras.Validate(); err != nil {
		return 0, err
	}
	if hlog.IsElapseFunction() {
		defer hlog.ElapseWithCtx(ctx, "hyper.shipping.Create",
			hlog.F(zap.String("uniLink", paras.UniLink.Str())),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Any("paras", paras), zap.Error(err)}
				}
				return nil
			},
		)()
	}
	id = NxtShippingID()
	exM := ex.MustEx(paras.Ex)
	shippingM := &ShipM{
		Template: hdb.Template{
			Ctrl: exM.Ctrl,
			Tag:  hjson.MustToBytes(exM.Tag),
			Meta: hjson.MustToBytes(exM.Meta),
		},
		ID:           id,
		UniLink:      paras.UniLink,
		CourierCode:  "",
		TrackingNo:   "",
		ShippedAt:    nil,
		LastSyncedAt: nil,
		Status:       StatusCreated,
		Address:      hjson.MustToBytes(paras.Address),
	}
	err = hdb.Create[ShipM](hyperplt.Tx(ctx), shippingM)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func Prepared(
	ctx context.Context,
	uniLink collar.Link,
	courierCode CourierCode,
	trackingNo string,
	meta ex.Meta,
) (err error) {
	if uniLink == "" {
		return errors.New("uni_link is required")
	}
	if courierCode == "" {
		return errors.New("courier_code is required")
	}
	if trackingNo == "" {
		return errors.New("tracking_no is required")
	}
	if hlog.IsElapseComponent() {
		defer hlog.ElapseWithCtx(ctx, "hyper.shipping.Prepared",
			hlog.F(zap.String("uni_link", uniLink.Str())),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{
						zap.String("uni_link", uniLink.Display()),
						zap.String("courierCode", courierCode),
						zap.String("trackingNo", trackingNo),
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
	mut := map[string]interface{}{
		"status":       StatusPickedUp,
		"courier_code": courierCode,
		"tracking_no":  trackingNo,
		"shipped_at":   gorm.Expr("CURRENT_TIMESTAMP"),
	}
	if len(meta) > 0 {
		if shipM.Meta == nil {
			mut["meta"] = hjson.MustToBytes(meta)
		} else {
			dbMeta := *hjson.MustFromBytes[ex.Meta](shipM.Meta)
			if dbMeta != nil {
				for k, v := range meta {
					dbMeta[k] = v
				}
			}
			mut["meta"] = hjson.MustToBytes(dbMeta)
		}
	}
	err = hdb.Update[ShipM](tx, mut, "id = ?", shipM.ID)
	if err != nil {
		return err
	}
	return nil
}

func Completed(
	ctx context.Context,
	uniLink collar.Link,
	meta ex.Meta,
) (err error) {
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
	mut := map[string]interface{}{
		"status": StatusDelivered,
	}
	if len(meta) > 0 {
		if shipM.Meta == nil {
			mut["meta"] = hjson.MustToBytes(meta)
		} else {
			dbMeta := *hjson.MustFromBytes[ex.Meta](shipM.Meta)
			if dbMeta != nil {
				for k, v := range meta {
					dbMeta[k] = v
				}
			}
			mut["meta"] = hjson.MustToBytes(dbMeta)
		}
	}
	err = hdb.Update[ShipM](tx, mut, "id = ?", shipM.ID)
	if err != nil {
		return err
	}
	return nil
}
