package brand

import (
	"github.com/hootuu/helix/storage/hpg"
	"gorm.io/datatypes"
)

type ID = string

type BrM struct {
	hpg.Basic
	ID          ID             `gorm:"column:id;primaryKey;size:32;"`
	Name        string         `gorm:"column:name;index;size:100;"`
	Intro       string         `gorm:"column:intro;size:1000;"`
	Description string         `gorm:"column:description;type:text;"`
	Media       datatypes.JSON `gorm:"column:media;type:jsonb;"`
}

func (m *BrM) TableName() string {
	return "hyper_brand"
}
