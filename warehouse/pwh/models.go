package pwh

import (
	"github.com/hootuu/helix/storage/hpg"
	"gorm.io/datatypes"
)

type PhysicalWhM struct {
	hpg.Basic
	ID     uint64 `gorm:"column:id;primaryKey;"`
	Collar string `gorm:"column:collar;uniqueIndex;size:64;"`
	Memo   string `gorm:"column:memo;size:100;"`
}

func (m *PhysicalWhM) TableName() string {
	return "hyper_warehouse_pwh"
}

type PhysicalSkuM struct {
	hpg.Basic
	PWH       uint64      `gorm:"column:pwh;uniqueIndex:uk_pwh_sku;"`
	SKU       uint64      `gorm:"column:sku;uniqueIndex:uk_pwh_sku;"`
	Available uint64      `gorm:"column:available;"`
	Locked    uint64      `gorm:"column:locked;"`
	Version   hpg.Version `gorm:"column:version;"`
}

func (m *PhysicalSkuM) TableName() string {
	return "hyper_warehouse_pwh_sku"
}

type PhysicalInOutM struct {
	hpg.Basic
	PWH       uint64         `gorm:"column:pwh;index;"`
	SKU       uint64         `gorm:"column:sku;index;"`
	Direction Direction      `gorm:"column:direction;"`
	Quantity  uint64         `gorm:"column:quantity;"`
	Price     uint64         `gorm:"column:price;"`
	Meta      datatypes.JSON `gorm:"column:meta;type:jsonb;"`
}

func (m *PhysicalInOutM) TableName() string {
	return "hyper_warehouse_pwh_in_out"
}
