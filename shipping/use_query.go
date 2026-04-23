package shipping

import (
	"context"
	"errors"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hypes/ex"
	"github.com/hootuu/hyper/hyperplt"
	"gorm.io/gorm"
)

func GetByOrderID(ctx context.Context, bizID string) (*ShippingInfo, error) {
	InitIfNeeded()
	if bizID == "" {
		return nil, errors.New("biz_id is required")
	}

	tx := hyperplt.Tx(ctx)
	shipM, err := hdb.Get[ShipM](tx, "biz_id = ?", bizID)
	if err != nil {
		return nil, err
	}
	if shipM == nil {
		return nil, nil
	}

	info := &ShippingInfo{
		Shipping: toShipping(shipM),
	}

	// Prefer package data from the new table.
	pkgArr, err := hdb.Find[ShipPkgM](func() *gorm.DB {
		return tx.Model(&ShipPkgM{}).
			Where("shipping_id = ?", shipM.ID).
			Order("package_seq ASC")
	})
	if err != nil {
		return nil, err
	}
	if len(pkgArr) > 0 {
		packages := make([]*ShippingPackage, 0, len(pkgArr))
		for _, pkg := range pkgArr {
			packages = append(packages, &ShippingPackage{
				ID:          pkg.ID,
				ShippingID:  pkg.ShippingID,
				PackageSeq:  pkg.PackageSeq,
				CourierCode: pkg.CourierCode,
				TrackingNo:  pkg.TrackingNo,
				IsPrimary:   pkg.IsPrimary,
			})
		}
		info.Packages = packages
		return info, nil
	}

	// Fallback for legacy data: map old main logistics fields into one package.
	if shipM.CourierCode != "" || shipM.TrackingNo != "" {
		info.Packages = []*ShippingPackage{
			{
				ID:          shipM.ID,
				ShippingID:  shipM.ID,
				PackageSeq:  1,
				CourierCode: shipM.CourierCode,
				TrackingNo:  shipM.TrackingNo,
				IsPrimary:   true,
			},
		}
	}
	return info, nil
}

func toShipping(shipM *ShipM) *Shipping {
	if shipM == nil {
		return nil
	}
	var addr *Address
	if len(shipM.Address) > 0 {
		addr = hjson.MustFromBytes[Address](shipM.Address)
	}
	return &Shipping{
		ID:          shipM.ID,
		Address:     addr,
		CourierCode: shipM.CourierCode,
		TrackingNo:  shipM.TrackingNo,
		ShippedAt:   shipM.SubmittedTime,
		DeliveredAt: shipM.CompletedTime,
		Status:      shipM.Status,
		Ex:          ex.WithBytes(shipM.Ctrl, shipM.Tag, shipM.Meta),
	}
}
