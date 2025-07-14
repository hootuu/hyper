package product

import (
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyper/brand"
	"github.com/hootuu/hyper/category"
	"github.com/hootuu/hyper/spec"
	"gorm.io/datatypes"
)

type SpuM struct {
	hdb.Basic
	Collar    collar.Collar  `gorm:"column:collar;index;size:64;"`
	ID        SpuID          `gorm:"column:id;primaryKey;size:32;"`
	Category  category.ID    `gorm:"column:category;index;"`
	Name      string         `gorm:"column:name;size:100;"`
	Intro     string         `gorm:"column:intro;size:1000;"`
	Brand     brand.ID       `gorm:"column:brand;size:32;"`
	Version   hdb.Version    `gorm:"column:version;"`
	MainMedia datatypes.JSON `gorm:"column:main_media;type:jsonb;"` //media.More
	MoreMedia datatypes.JSON `gorm:"column:more_media;type:jsonb;"` //media.Dict
}

func (m *SpuM) TableName() string {
	return "hyper_product_spu"
}

type SpuSpecM struct {
	hdb.Basic
	Spu  SpuID   `gorm:"column:spu;uniqueIndex:uk_spu_spec;size:32;"`
	Spec spec.ID `gorm:"column:spec;uniqueIndex:uk_spu_spec;"`
	Seq  int     `gorm:"column:seq;"`
}

func (m *SpuSpecM) TableName() string {
	return "hyper_product_spu_spec"
}

type SpuSpecOptM struct {
	hdb.Basic
	ID    spec.OptID     `gorm:"column:id;primaryKey;autoIncrement:false;"`
	Spu   SpuID          `gorm:"column:spu;uniqueIndex:uk_spu_spec;size:32;"`
	Spec  spec.ID        `gorm:"column:spec;uniqueIndex:uk_spu_spec;"`
	Label string         `gorm:"column:name;size:100;"`
	Media datatypes.JSON `gorm:"column:media;type:jsonb;"` //media.More
	Seq   int            `gorm:"column:seq;"`
}

func (m *SpuSpecOptM) TableName() string {
	return "hyper_product_spu_spec_opt"
}

type SkuM struct {
	hdb.Basic
	ID  SkuID `gorm:"column:id;primaryKey;size:32;"`
	Spu SpuID `gorm:"column:spu;index;size:32;"`
}

func (m *SkuM) TableName() string {
	return "hyper_product_sku"
}

type SkuSpecM struct {
	hdb.Basic
	Sku     SkuID      `gorm:"column:sku;uniqueIndex:uk_sku_spec;size:32;"`
	SpecOpt spec.OptID `gorm:"column:spec_opt;uniqueIndex:uk_sku_spec;"`
}

func (m *SkuSpecM) TableName() string {
	return "hyper_product_sku_spec"
}
