package pwh

import (
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyper/hiprod/prod"
	"gorm.io/datatypes"
)

type PhysicalWhM struct {
	hdb.Template
	ID   ID        `gorm:"column:id;primaryKey;autoIncrement:false;"`
	Link collar.ID `gorm:"column:link;uniqueIndex;size:64;"`
	Memo string    `gorm:"column:memo;size:100;"`
}

func (m *PhysicalWhM) TableName() string {
	return "hyper_prod_pwh"
}

type PhysicalSkuM struct {
	hdb.Basic
	PWH       ID          `gorm:"column:pwh;uniqueIndex:uk_pwh_sku;autoIncrement:false;"`
	SKU       prod.SkuID  `gorm:"column:sku;uniqueIndex:uk_pwh_sku;autoIncrement:false;"`
	Available uint64      `gorm:"column:available;"`
	Locked    uint64      `gorm:"column:locked;"`
	Version   hdb.Version `gorm:"column:version;"`
}

func (m *PhysicalSkuM) TableName() string {
	return "hyper_prod_pwh_sku"
}

type PhysicalInOutM struct {
	hdb.Basic
	PWH       ID             `gorm:"column:pwh;index;autoIncrement:false;"`
	SKU       prod.SkuID     `gorm:"column:sku;index;autoIncrement:false;"`
	Direction Direction      `gorm:"column:direction;"`
	Quantity  uint64         `gorm:"column:quantity;"`
	Price     uint64         `gorm:"column:price;"`
	Meta      datatypes.JSON `gorm:"column:meta;type:jsonb;"`
}

func (m *PhysicalInOutM) TableName() string {
	return "hyper_prod_pwh_in_out"
}

type PhysicalLockUnlockM struct {
	hdb.Basic
	PWH       ID             `gorm:"column:pwh;index;autoIncrement:false;"`
	SKU       prod.SkuID     `gorm:"column:sku;index;autoIncrement:false;"`
	Direction Direction      `gorm:"column:direction;"`
	Quantity  uint64         `gorm:"column:quantity;"`
	Meta      datatypes.JSON `gorm:"column:meta;type:jsonb;"`
}

func (m *PhysicalLockUnlockM) TableName() string {
	return "hyper_prod_pwh_lock_unlock"
}
