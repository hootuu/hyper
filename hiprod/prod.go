package hiprod

import (
	"context"
	"errors"
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
	Biz       prod.Biz    `json:"biz"`      //optional
	Category  category.ID `json:"category"` //optional
	Brand     brand.ID    `json:"brand"`    //optional
	Name      string      `json:"name"`
	Intro     string      `json:"intro"` //optional
	Media     media.Dict  `json:"media"` //optional
	Price     uint64      `json:"price"`
	Cost      uint64      `json:"cost"` //optional
	Inventory uint64      `json:"inventory"`
	PwhID     pwh.ID      `json:"pwh_id"` //optional
	VwhID     vwh.ID      `json:"vwh_id"` //optional
	Ctrl      ctrl.Ctrl   `json:"ctrl"`
	Tag       tag.Tag     `json:"tag"`
	Meta      dict.Dict   `json:"meta"`
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
			Biz:      paras.Biz,
			Category: paras.Category,
			Name:     paras.Name,
			Intro:    paras.Intro,
			Brand:    paras.Brand,
			Media:    hjson.MustToBytes(paras.Media),
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
