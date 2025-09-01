package hiprod

import (
	"context"
	"errors"
	"fmt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/ctrl"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/data/tag"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/media"
	"github.com/hootuu/hyper/brand"
	"github.com/hootuu/hyper/category"
	"github.com/hootuu/hyper/hiprod/prod"
	"github.com/hootuu/hyper/hiprod/pwh"
	"github.com/hootuu/hyper/hiprod/vwh"
	"github.com/hootuu/hyper/hyperplt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProdCreateParas struct {
	ID        prod.ID     `json:"id"`
	Biz       prod.Biz    `json:"biz"`      //optional
	Category  category.ID `json:"category"` //optional
	Brand     brand.ID    `json:"brand"`    //optional
	Name      string      `json:"name"`     //require
	Intro     string      `json:"intro"`    //optional
	Media     media.Dict  `json:"media"`    //optional
	Price     uint64      `json:"price"`    //*
	Cost      uint64      `json:"cost"`     //optional
	Available bool        `json:"available"`
	Inventory uint64      `json:"inventory"` //*
	PwhID     pwh.ID      `json:"pwh_id"`    //optional
	VwhID     vwh.ID      `json:"vwh_id"`    //optional
	Ctrl      ctrl.Ctrl   `json:"ctrl"`      //optional
	Tag       tag.Tag     `json:"tag"`       //optional
	Meta      dict.Dict   `json:"meta"`      //optional
}

func (p *ProdCreateParas) validate() error {
	if p.Biz == "" {
		p.Biz = prod.UniBiz
	}
	if p.Name == "" {
		return errors.New("name is required")
	}
	if p.Inventory == 0 {
		return errors.New("inventory is required")
	}
	if p.PwhID == 0 {
		p.PwhID = pwh.UniPwhID()
	}
	if p.VwhID == 0 {
		p.VwhID = vwh.UniVwhID()
	}
	return nil
}

func CreateProduct(ctx context.Context, paras *ProdCreateParas) (skuID prod.SkuID, err error) {
	if paras == nil {
		return 0, errors.New("paras is required")
	}
	if err := paras.validate(); err != nil {
		return 0, err
	}
	defer hlog.ElapseWithCtx(ctx, "hyper.prod.CreateProduct",
		hlog.F(zap.String("paras.name", paras.Name)),
		func() []zap.Field {
			if err != nil {
				return []zap.Field{zap.Any("paras", paras), zap.Error(err)}
			}
			return []zap.Field{zap.Uint64("skuID", skuID)}
		})()
	tx := hyperplt.Tx(ctx)
	err = hdb.Tx(tx, func(tx *gorm.DB) error {
		innerCtx := hdb.TxCtx(tx, ctx)
		var err error
		spuM, err := prod.CreateSpu(innerCtx, &prod.SpuM{
			Template: hdb.Template{
				Ctrl: paras.Ctrl,
				Tag:  hjson.MustToBytes(paras.Tag),
				Meta: hjson.MustToBytes(paras.Meta),
			},
			Biz:       paras.Biz,
			Category:  paras.Category,
			Name:      paras.Name,
			Intro:     paras.Intro,
			Brand:     paras.Brand,
			Media:     hjson.MustToBytes(paras.Media),
			Cost:      paras.Cost,
			Available: paras.Available,
		})
		if err != nil {
			return err
		}
		skuID, err = prod.CreateSku(innerCtx, &prod.SkuSpecSetting{
			Spu:   spuM.ID,
			Specs: nil,
		})
		err = pwh.Into(innerCtx, pwh.IntoOutParas{
			PwhID:    paras.PwhID,
			SkuID:    skuID,
			Quantity: paras.Inventory,
			Price:    paras.Cost,
			Meta:     paras.Meta,
		})
		if err != nil {
			return err
		}
		err = vwh.SetSku(innerCtx, vwh.SetSkuParas{
			Vwh:       paras.VwhID,
			Sku:       skuID,
			Pwh:       paras.PwhID,
			Price:     paras.Price,
			Inventory: paras.Inventory,
		})
		return nil
	})

	if err != nil {
		return 0, err
	}
	return skuID, nil
}

func CreateProductByPwh(ctx context.Context, paras *ProdCreateParas) (skuID prod.SkuID, err error) {
	if paras == nil {
		return 0, errors.New("paras is required")
	}
	if err := paras.validate(); err != nil {
		return 0, err
	}
	defer hlog.ElapseWithCtx(ctx, "hyper.prod.CreateProduct",
		hlog.F(zap.String("paras.name", paras.Name)),
		func() []zap.Field {
			if err != nil {
				return []zap.Field{zap.Any("paras", paras), zap.Error(err)}
			}
			return []zap.Field{zap.Uint64("skuID", skuID)}
		})()
	tx := hyperplt.Tx(ctx)
	err = hdb.Tx(tx, func(tx *gorm.DB) error {
		innerCtx := hdb.TxCtx(tx, ctx)
		var err error
		spuM, err := prod.CreateSpu(innerCtx, &prod.SpuM{
			Template: hdb.Template{
				Ctrl: paras.Ctrl,
				Tag:  hjson.MustToBytes(paras.Tag),
				Meta: hjson.MustToBytes(paras.Meta),
			},
			Biz:       paras.Biz,
			Category:  paras.Category,
			Name:      paras.Name,
			Intro:     paras.Intro,
			Brand:     paras.Brand,
			Media:     hjson.MustToBytes(paras.Media),
			Cost:      paras.Cost,
			Price:     paras.Price,
			Available: paras.Available,
		})
		if err != nil {
			return err
		}
		skuID, err = prod.CreateSku(innerCtx, &prod.SkuSpecSetting{
			Spu:   spuM.ID,
			Specs: nil,
		})
		err = pwh.Into(innerCtx, pwh.IntoOutParas{
			PwhID:    paras.PwhID,
			SkuID:    skuID,
			Quantity: paras.Inventory,
			Price:    paras.Cost,
		})
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return 0, err
	}
	return skuID, nil
}

type ProdUpdateParas struct {
	SpuID prod.SpuID `json:"spu_id"`
	Price uint64     `json:"price"`
	Cost  uint64     `json:"cost"`
}

func (p *ProdUpdateParas) validate() error {
	if p.SpuID == 0 {
		return errors.New("require SpuID")
	}
	if p.Price == 0 && p.Cost == 0 {
		return errors.New("require Price or Cost")
	}
	return nil
}

func UpdateProductPrice(ctx context.Context, paras *ProdUpdateParas) (err error) {
	if paras == nil {
		return errors.New("paras is required")
	}
	if err := paras.validate(); err != nil {
		return err
	}
	defer hlog.ElapseWithCtx(ctx, "hyper.prod.UpdateProductPrice",
		hlog.F(zap.Any("paras", paras)),
		func() []zap.Field {
			if err != nil {
				return []zap.Field{zap.Any("paras", paras), zap.Error(err)}
			}
			return nil
		})()

	tx := hyperplt.Tx(ctx)
	updateM := map[string]interface{}{}
	if paras.Price > 0 {
		updateM["price"] = paras.Price
	}
	if paras.Cost > 0 {
		updateM["cost"] = paras.Cost
	}
	err = hdb.Update[prod.SpuM](tx, updateM, "id = ?", paras.SpuID)
	return err
}

func SetAvailable(ctx context.Context, spuID prod.SpuID, available bool) (err error) {
	if spuID == 0 {
		return errors.New("require spuID")
	}
	tx := hyperplt.Tx(ctx)
	err = hdb.Update[prod.SpuM](tx, map[string]interface{}{
		"available": available,
	}, "id = ?", spuID)
	return err
}

type SkuUnpPublishParas struct {
	Vwh  vwh.ID     `json:"vwh"`
	Sku  prod.SkuID `json:"sku"`
	Pwh  pwh.ID     `json:"pwh"`
	Meta dict.Dict  `json:"meta"`
}

func (p SkuUnpPublishParas) Validate() error {
	if p.Vwh == 0 {
		return fmt.Errorf("require vwh")
	}
	if p.Pwh == 0 {
		return fmt.Errorf("require pwh")
	}
	if p.Sku == 0 {
		return fmt.Errorf("require sku")
	}
	return nil
}

func SkuUnpPublish(ctx context.Context, paras SkuUnpPublishParas) error {
	err := paras.Validate()
	if err != nil {
		return err
	}
	tx := hyperplt.Tx(ctx)
	vwhSkuM, err := hdb.Get[vwh.VirtualWhSkuM](tx, "vwh = ? AND sku = ? AND pwh = ?",
		paras.Vwh, paras.Sku, paras.Pwh)
	if err != nil {
		return err
	}
	if vwhSkuM == nil {
		return fmt.Errorf("vwh sku not found: vwh=%d, sku=%d, pwh=%d", paras.Vwh, paras.Sku, paras.Pwh)
	}
	vwhSkuExtM, err := hdb.Get[vwh.VirtualWhSkuExtM](tx, "pwh = ? AND link = ?", paras.Pwh, vwh.BuildVwhExtLink(paras.Vwh, paras.Sku))
	if err != nil {
		return err
	}

	err = hdb.Tx(tx, func(innerTx *gorm.DB) error {
		quantity := vwhSkuM.Inventory
		if quantity > 0 {
			err = pwh.Unlock(hdb.TxCtx(innerTx), pwh.LockUnlockParas{
				PwhID:    paras.Pwh,
				SkuID:    paras.Sku,
				Quantity: quantity,
				Meta:     paras.Meta,
			})
			if err != nil {
				return err
			}
		}
		err = hdb.Update[vwh.VirtualWhSkuM](innerTx, map[string]any{
			"inventory": 0,
			"version":   gorm.Expr("version + 1"),
		}, "auto_id = ?", vwhSkuM.AutoID)
		if err != nil {
			return err
		}
		if vwhSkuExtM != nil {
			err = hdb.Update[vwh.VirtualWhSkuExtM](innerTx, map[string]any{
				"available": false,
			}, "auto_id = ?", vwhSkuExtM.AutoID)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
