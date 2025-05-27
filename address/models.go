package address

import (
	"github.com/hootuu/helix/components/htree"
	"github.com/hootuu/helix/storage/hpg"
	"github.com/hootuu/hyper/address/maps"
)

type DistrictM struct {
	hpg.Basic
	ID   htree.ID `gorm:"column:id;primaryKey;"`
	Map  maps.Map `gorm:"column:map;uniqueIndex:uk_map_code;"`
	Code string   `gorm:"column:code;uniqueIndex:uk_map_code;not null;size:20;"`
	Name string   `gorm:"column:name;not null;size:32;"`
}

func (m *DistrictM) TableName() string {
	return "hyper_address_district"
}
