package brand

import (
	"github.com/hootuu/helix/storage/hdb"
	"gorm.io/datatypes"
)

type ID = uint64

type BrM struct {
	hdb.Basic
	ID          ID             `gorm:"column:id;primaryKey;autoIncrement:false;"`
	Name        string         `gorm:"column:name;index;size:100;"`
	Intro       string         `gorm:"column:intro;size:1000;"`
	Description string         `gorm:"column:description;type:text;"`
	Media       datatypes.JSON `gorm:"column:media;type:jsonb;"`
}

func (m *BrM) TableName() string {
	return "hyper_brand"
}
