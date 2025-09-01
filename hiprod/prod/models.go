package prod

import (
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyper/brand"
	"github.com/hootuu/hyper/category"
	"gorm.io/datatypes"
)

type SpuM struct {
	hdb.Template
	Biz       Biz            `gorm:"column:biz;index;size:16;"`
	ID        SpuID          `gorm:"column:id;primaryKey;autoIncrement:false;"`
	Category  category.ID    `gorm:"column:category;index;autoIncrement:false;"`
	Name      string         `gorm:"column:name;size:100;"`
	Intro     string         `gorm:"column:intro;size:618;"`
	Brand     brand.ID       `gorm:"column:brand;size:32;"`
	Version   hdb.Version    `gorm:"column:version;"`
	Media     datatypes.JSON `gorm:"column:media;type:jsonb;"` //media.Dict
	Cost      uint64         `gorm:"column:cost;size:32;"`
	Price     uint64         `gorm:"column:price;size:32;"`
	Available bool           `gorm:"column:available;"`
}

func (m *SpuM) TableName() string {
	return "hyper_prod_spu"
}

type SpuSpecM struct {
	hdb.Basic
	Spu  SpuID  `gorm:"column:spu;uniqueIndex:uk_spu_spec;size:32;"`
	Spec SpecID `gorm:"column:spec;uniqueIndex:uk_spu_spec;"`
	Seq  int    `gorm:"column:seq;"`
}

func (m *SpuSpecM) TableName() string {
	return "hyper_prod_spu_spec"
}

type SpuSpecOptM struct {
	hdb.Basic
	ID    SpecOptID      `gorm:"column:id;primaryKey;autoIncrement:false;"`
	Spu   SpuID          `gorm:"column:spu;uniqueIndex:uk_spu_spec;"`
	Spec  SpecID         `gorm:"column:spec;uniqueIndex:uk_spu_spec;autoIncrement:false;"`
	Label string         `gorm:"column:name;size:100;"`
	Media datatypes.JSON `gorm:"column:media;type:jsonb;"` //media.More
	Seq   int            `gorm:"column:seq;"`
}

func (m *SpuSpecOptM) TableName() string {
	return "hyper_prod_spu_spec_opt"
}

type SkuM struct {
	hdb.Basic
	ID  SkuID `gorm:"column:id;primaryKey;autoIncrement:false;"`
	Spu SpuID `gorm:"column:spu;index;autoIncrement:false;"`
}

func (m *SkuM) TableName() string {
	return "hyper_product_sku"
}

type SkuSpecM struct {
	hdb.Basic
	Sku     SkuID     `gorm:"column:sku;uniqueIndex:uk_sku_spec;size:32;autoIncrement:false;"`
	SpecOpt SpecOptID `gorm:"column:spec_opt;uniqueIndex:uk_sku_spec;autoIncrement:false;"`
}

func (m *SkuSpecM) TableName() string {
	return "hyper_prod_sku_spec"
}
