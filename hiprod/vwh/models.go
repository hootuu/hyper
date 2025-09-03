package vwh

import (
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyper/hiprod/prod"
	"github.com/hootuu/hyper/hiprod/pwh"
	"gorm.io/datatypes"
)

type VirtualWhM struct {
	hdb.Template
	ID   ID        `gorm:"column:id;primaryKey;autoIncrement:false;"`
	Link collar.ID `gorm:"column:link;uniqueIndex;size:64;"`
	Memo string    `gorm:"column:memo;size:100;"`
}

func (m *VirtualWhM) TableName() string {
	return "hyper_prod_vwh"
}

type VirtualWhSrcM struct {
	hdb.Basic
	Vwh           ID             `gorm:"column:vwh;uniqueIndex:uk_vwh_pwh;autoIncrement:false;"`
	Pwh           pwh.ID         `gorm:"column:pwh;uniqueIndex:uk_vwh_pwh;autoIncrement:false;"`
	PricingType   string         `gorm:"column:pricing_type;size:64;"`
	InventoryType string         `gorm:"column:inventory_type;size:64;"`
	Pricing       datatypes.JSON `gorm:"column:pricing;type:jsonb;"`
	Inventory     datatypes.JSON `gorm:"column:inventory;type:jsonb;"`
}

func (m *VirtualWhSrcM) TableName() string {
	return "hyper_prod_vwh_src"
}

type VirtualWhSkuM struct {
	hdb.Basic
	Vwh          ID          `gorm:"column:vwh;uniqueIndex:uk_vwh_pwh_sku;autoIncrement:false;"`
	Sku          prod.SkuID  `gorm:"column:sku;uniqueIndex:uk_vwh_pwh_sku;autoIncrement:false;"`
	Pwh          pwh.ID      `gorm:"column:pwh;uniqueIndex:uk_vwh_pwh_sku;autoIncrement:false;"`
	UseInventory bool        `gorm:"column:use_inventory;size:1;"`
	Price        uint64      `gorm:"column:price;"`
	Inventory    uint64      `gorm:"column:inventory;"`
	Version      hdb.Version `gorm:"column:version;"`
}

func (m *VirtualWhSkuM) TableName() string {
	return "hyper_prod_vwh_sku"
}

type VirtualWhSkuExtM struct {
	hdb.Template
	Vwh       ID         `gorm:"column:vwh;uniqueIndex:uk_vwh_pwh_sku;autoIncrement:false;"`
	Sku       prod.SkuID `gorm:"column:sku;uniqueIndex:uk_vwh_pwh_sku;autoIncrement:false;"`
	Pwh       pwh.ID     `gorm:"column:pwh;uniqueIndex:uk_vwh_pwh_sku;autoIncrement:false;"`
	Link      collar.ID  `gorm:"column:link;size:64;uniqueIndex:uk_link;"`
	Channel   uint64     `gorm:"column:channel;"`
	Available bool       `gorm:"column:available;size:1;"`
	Sort      uint64     `gorm:"column:sort;"`
	Sales     uint64     `gorm:"column:sales;"`
}

func (m *VirtualWhSkuExtM) TableName() string {
	return "hyper_prod_vwh_sku_ext"
}
