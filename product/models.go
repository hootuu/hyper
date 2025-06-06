package product

import (
	"github.com/hootuu/helix/storage/hpg"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyper/brand"
	"github.com/hootuu/hyper/category"
	"gorm.io/datatypes"
)

type SpuM struct {
	hpg.Basic
	Collar    collar.Collar  `gorm:"column:collar;index;size:64;"`
	ID        SpuID          `gorm:"column:id;primaryKey;size:32;"`
	Frontend  category.ID    `gorm:"column:frontend;index;"`
	Backend   category.ID    `gorm:"column:backend;index;"`
	Name      string         `gorm:"column:name;size:100;"`
	Intro     string         `gorm:"column:intro;size:1000;"`
	Brand     brand.ID       `gorm:"column:brand;size:32;"`
	Version   hpg.Version    `gorm:"column:version;"`
	MainMedia datatypes.JSON `gorm:"column:main_media;type:jsonb;"` //media.More
	MoreMedia datatypes.JSON `gorm:"column:more_media;type:jsonb;"` //media.Dict
}

func (m *SpuM) TableName() string {
	return "hyper_product_spu"
}

type SpuSpecM struct {
	hpg.Basic
	Spu   SpuID  `gorm:"column:spu;uniqueIndex:uk_spu_spec;size:32;"`
	Spec  SpecID `gorm:"column:spec;uniqueIndex:uk_spu_spec;"`
	Name  string `gorm:"column:name;size:50;"`
	Intro string `gorm:"column:intro;size:300;"`
	Key   bool   `gorm:"column:is_key;"`
}

func (m *SpuSpecM) TableName() string {
	return "hyper_product_spu_spec"
}

type SkuM struct {
	hpg.Basic
	ID  SkuID `gorm:"column:id;primaryKey;size:32;"`
	Spu SpuID `gorm:"column:spu;index;size:32;"`
}

func (m *SkuM) TableName() string {
	return "hyper_product_sku"
}

type SkuSpecM struct {
	hpg.Basic
	Sku   SkuID          `gorm:"column:sku;uniqueIndex:uk_sku_spec;size:32;"`
	Spec  SpecID         `gorm:"column:spec;uniqueIndex:uk_sku_spec;"`
	Label string         `gorm:"column:name;size:100;"`
	Media datatypes.JSON `gorm:"column:media;type:jsonb;"` //media.More
}

func (m *SkuSpecM) TableName() string {
	return "hyper_product_sku_spec"
}
