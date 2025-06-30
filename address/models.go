package address

import (
	"github.com/hootuu/helix/components/htree"
	"github.com/hootuu/helix/components/sattva"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyper/address/maps"
)

type RegionM struct {
	hdb.Basic
	ID      htree.ID `gorm:"column:id;primaryKey;"`
	Map     maps.Map `gorm:"column:map;uniqueIndex:uk_map_code_name;"`
	Code    string   `gorm:"column:code;uniqueIndex:uk_map_code_name;not null;size:20;"`
	Name    string   `gorm:"column:name;uniqueIndex:uk_map_code_name;not null;size:32;"`
	Address string   `gorm:"column:address;not null;size:200;"`
}

func (m *RegionM) TableName() string {
	return "hyper_address_region"
}

type AddrM struct {
	hdb.Basic
	ID       string                `gorm:"column:id;primaryKey;"`
	Owner    sattva.Identification `gorm:"column:owner;index;not null;size:32;"`
	CName    string                `gorm:"column:cname;not null;size:32;"`
	CMobi    string                `gorm:"column:cmobi;not null;size:32;"`
	Default  bool                  `gorm:"column:is_default;"`
	Region   RegionID              `gorm:"column:region;index;"`
	Addr     string                `gorm:"column:addr;not null;size:100;"`
	FullAddr string                `gorm:"column:full_addr;not null;size:200;"`
	LocX     float64               `gorm:"column:loc_x;type:decimal(9,6);"`
	LocY     float64               `gorm:"column:loc_y;type:decimal(9,6);"`
	Usage    int64                 `gorm:"column:usage;"`
}

func (m *AddrM) TableName() string {
	return "hyper_address_addr"
}

func (m *AddrM) ToAddress() *Address {
	return &Address{
		ID:       m.ID,
		Owner:    m.Owner,
		Region:   m.Region,
		Address:  m.Addr,
		FullAddr: m.FullAddr,
		Contact: Contact{
			Name: m.CName,
			Mobi: m.CMobi,
		},
		Default: m.Default,
		Location: Location{
			Lon: m.LocX,
			Lat: m.LocY,
		},
	}
}
