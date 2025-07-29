package hiprod

import (
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/hypes/media"
	"github.com/hootuu/hyper/brand"
	"github.com/hootuu/hyper/category"
	"github.com/hootuu/hyper/hiprod/prod"
	"github.com/hootuu/hyper/hiprod/pwh"
	"github.com/hootuu/hyper/hiprod/vwh"
)

type Product struct {
	ID        prod.ID     `json:"id"`
	SkuID     prod.SkuID  `json:"sku_id"`
	SpuID     prod.SpuID  `json:"spu_id"`
	VwhID     vwh.ID      `json:"vwh_id"`
	PwhID     pwh.ID      `json:"pwh_id"`
	Biz       prod.Biz    `json:"biz"`
	Category  category.ID `json:"category"`
	Name      string      `json:"name"`
	Intro     string      `json:"intro"`
	Brand     brand.ID    `json:"brand"`
	Media     media.Dict  `json:"media"`
	Price     uint64      `gorm:"column:price;"`
	Inventory uint64      `gorm:"column:inventory;"`
}

func (prod *Product) GetDigest() dict.Dict {
	return dict.New(map[string]any{
		"sku_id":    prod.SkuID,
		"spu_id":    prod.SpuID,
		"vwh_id":    prod.VwhID,
		"pwh_id":    prod.PwhID,
		"biz":       prod.Biz,
		"category":  prod.Category,
		"name":      prod.Name,
		"intro":     prod.Intro,
		"brand":     prod.Brand,
		"price":     prod.Price,
		"inventory": prod.Inventory,
		"media":     prod.Media,
	})
}
