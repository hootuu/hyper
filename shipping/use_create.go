package shipping

import (
	"context"
	"errors"
	"time"

	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/ex"
	"github.com/hootuu/hyper/hyperplt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CreateParas struct {
	Idem    string        `json:"idem"`
	BizCode string        `json:"biz_code"`
	BizID   string        `json:"biz_id"`
	Address *Address      `json:"address"`
	Ex      *ex.Ex        `json:"ex"`
	Timeout time.Duration `json:"timeout"`
}

func (paras *CreateParas) Validate() error {
	if paras.Idem == "" {
		return errors.New("idem is required")
	}
	if paras.BizCode == "" {
		return errors.New("biz_code is required")
	}
	if paras.BizID == "" {
		return errors.New("biz_id is required")
	}
	if paras.Address == nil {
		return errors.New("address is required")
	}
	if err := paras.Address.Validate(); err != nil {
		return err
	}
	if paras.Timeout == 0 {
		paras.Timeout = 7 * 24 * time.Hour
	}
	return nil
}

func Create(ctx context.Context, paras *CreateParas) (id ID, err error) {
	InitIfNeeded()
	if paras == nil {
		return 0, errors.New("paras is required")
	}
	if err := paras.Validate(); err != nil {
		return 0, err
	}
	if hlog.IsElapseComponent() {
		defer hlog.ElapseWithCtx(ctx, "hyper.shipping.Create",
			hlog.F(zap.String("bizCode", paras.BizCode)),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Any("paras", paras), zap.Error(err)}
				}
				return nil
			},
		)()
	}
	if err := hyperplt.Idem().MustCheck(paras.Idem); err != nil {
		return 0, err
	}

	id = nxtShippingID()
	exM := ex.MustEx(paras.Ex)
	shippingM := &ShipM{
		Template: hdb.Template{
			Ctrl: exM.Ctrl,
			Tag:  hjson.MustToBytes(exM.Tag),
			Meta: hjson.MustToBytes(exM.Meta),
		},
		ID:               id,
		BizCode:          paras.BizCode,
		BizID:            paras.BizID,
		Status:           Initialized,
		Address:          hjson.MustToBytes(paras.Address),
		Timeout:          paras.Timeout,
		TimeoutCompleted: false,
	}
	err = hdb.Create[ShipM](hyperplt.Tx(ctx), shippingM)
	if err != nil {
		return 0, err
	}
	return id, nil
}

type CreatePackageItem struct {
	CourierCode CourierCode `json:"courier_code"`
	TrackingNo  string      `json:"tracking_no"`
}

func (item *CreatePackageItem) Validate() error {
	if item == nil {
		return errors.New("package item is required")
	}
	if item.CourierCode == "" {
		return errors.New("courier_code is required")
	}
	if item.TrackingNo == "" {
		return errors.New("tracking_no is required")
	}
	return nil
}

type CreateBatchParas struct {
	Idem     string               `json:"idem"`
	BizCode  string               `json:"biz_code"`
	BizID    string               `json:"biz_id"`
	Address  *Address             `json:"address"`
	Packages []*CreatePackageItem `json:"packages"`
	Ex       *ex.Ex               `json:"ex"`
	Timeout  time.Duration        `json:"timeout"`
}

func (paras *CreateBatchParas) Validate() error {
	if paras.Idem == "" {
		return errors.New("idem is required")
	}
	if paras.BizCode == "" {
		return errors.New("biz_code is required")
	}
	if paras.BizID == "" {
		return errors.New("biz_id is required")
	}
	if paras.Address == nil {
		return errors.New("address is required")
	}
	if err := paras.Address.Validate(); err != nil {
		return err
	}
	if len(paras.Packages) == 0 {
		return errors.New("packages is required")
	}
	for i := range paras.Packages {
		if err := paras.Packages[i].Validate(); err != nil {
			return err
		}
	}
	if paras.Timeout == 0 {
		paras.Timeout = 7 * 24 * time.Hour
	}
	return nil
}

type CreateBatchResult struct {
	ShippingID ID   `json:"shipping_id"`
	PackageIDs []ID `json:"package_ids"`
}

// CreateBatch creates one shipping order with multiple express packages.
// The first package is written back to hyper_shipping as the primary logistics info for compatibility.
func CreateBatch(ctx context.Context, paras *CreateBatchParas) (result *CreateBatchResult, err error) {
	InitIfNeeded()
	if paras == nil {
		return nil, errors.New("paras is required")
	}
	if err := paras.Validate(); err != nil {
		return nil, err
	}
	if hlog.IsElapseComponent() {
		defer hlog.ElapseWithCtx(ctx, "hyper.shipping.CreateBatch",
			hlog.F(zap.String("bizCode", paras.BizCode)),
			hlog.F(zap.String("bizId", paras.BizID)),
			hlog.F(zap.Any("Idem", paras.Idem)),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Any("paras", paras), zap.Error(err)}
				}
				return nil
			},
		)()
	}
	if err := hyperplt.Idem().MustCheck(paras.Idem); err != nil {
		return nil, err
	}
	result = &CreateBatchResult{
		PackageIDs: make([]ID, 0, len(paras.Packages)),
	}

	err = hdb.Tx(hyperplt.Tx(ctx), func(tx *gorm.DB) error {
		shippingID := nxtShippingID()
		exM := ex.MustEx(paras.Ex)
		primary := paras.Packages[0]
		submittedAt := time.Now()

		shippingM := &ShipM{
			Template: hdb.Template{
				Ctrl: exM.Ctrl,
				Tag:  hjson.MustToBytes(exM.Tag),
				Meta: hjson.MustToBytes(exM.Meta),
			},
			ID:               shippingID,
			BizCode:          paras.BizCode,
			BizID:            paras.BizID,
			CourierCode:      primary.CourierCode,
			TrackingNo:       primary.TrackingNo,
			Status:           Submitted,
			Address:          hjson.MustToBytes(paras.Address),
			Timeout:          paras.Timeout,
			TimeoutCompleted: false,
			SubmittedTime:    &submittedAt,
		}
		if err = hdb.Create[ShipM](tx, shippingM); err != nil {
			return err
		}
		result.ShippingID = shippingID

		pkgArr := make([]*ShipPkgM, 0, len(paras.Packages))
		for i, pkg := range paras.Packages {
			pkgID := nxtShippingID()
			result.PackageIDs = append(result.PackageIDs, pkgID)
			pkgArr = append(pkgArr, &ShipPkgM{
				Template: hdb.Template{
					Ctrl: exM.Ctrl,
					Tag:  hjson.MustToBytes(exM.Tag),
					Meta: hjson.MustToBytes(exM.Meta),
				},
				ID:          pkgID,
				ShippingID:  shippingID,
				BizCode:     paras.BizCode,
				BizID:       paras.BizID,
				PackageSeq:  i + 1,
				CourierCode: pkg.CourierCode,
				TrackingNo:  pkg.TrackingNo,
				IsPrimary:   i == 0,
			})
		}
		if err = hdb.MultiCreate[ShipPkgM](tx, pkgArr); err != nil {
			return err
		}
		mqPublishShippingAlter(&AlterPayload{
			ShippingID: shippingID,
			BizCode:    paras.BizCode,
			BizID:      paras.BizID,
			Src:        Initialized,
			Dst:        Submitted,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
