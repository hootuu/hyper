package vwh

import (
	"github.com/hootuu/helix/storage/hpg"
	"github.com/hootuu/hyper/sku"
	"github.com/hootuu/hyper/warehouse/pwh"
	"gorm.io/datatypes"
)

type VirtualWhM struct {
	hpg.Basic
	ID     ID             `gorm:"column:id;primaryKey;"`
	Collar string         `gorm:"column:collar;uniqueIndex;size:64;"`
	Meta   datatypes.JSON `gorm:"column:meta;type:jsonb;"`
}

func (m *VirtualWhM) TableName() string {
	return "hyper_warehouse_vwh"
}

type VirtualWhSrcM struct {
	hpg.Basic
	Vwh           ID             `gorm:"column:vwh;uniqueIndex:uk_vwh_pwh;"`
	Pwh           pwh.ID         `gorm:"column:pwh;uniqueIndex:uk_vwh_pwh;"`
	PricingType   string         `gorm:"column:pricing_type;size:64;"`
	InventoryType string         `gorm:"column:inventory_type;size:64;"`
	Pricing       datatypes.JSON `gorm:"column:pricing;type:jsonb;"`
	Inventory     datatypes.JSON `gorm:"column:inventory;type:jsonb;"`
}

func (m *VirtualWhSrcM) TableName() string {
	return "hyper_warehouse_vwh_src"
}

type VirtualWhSkuM struct {
	hpg.Basic
	Vwh       ID          `gorm:"column:vwh;uniqueIndex:uk_vwh_pwh_sku;"`
	Sku       sku.ID      `gorm:"column:sku;uniqueIndex:uk_vwh_pwh_sku;"`
	Pwh       pwh.ID      `gorm:"column:pwh;uniqueIndex:uk_vwh_pwh_sku;"`
	Price     uint64      `gorm:"column:price;"`
	Inventory uint64      `gorm:"column:inventory;"`
	Version   hpg.Version `gorm:"column:version;"`
}

func (m *VirtualWhSkuM) TableName() string {
	return "hyper_warehouse_vwh_sku"
}
