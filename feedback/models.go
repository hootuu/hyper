package feedback

import (
	"github.com/hootuu/helix/components/sattva"
	"github.com/hootuu/helix/storage/hdb"
	"gorm.io/datatypes"
)

type FbM struct {
	hdb.Basic
	ID          string                `gorm:"column:id;primaryKey;size:32;"`
	Person      sattva.Identification `gorm:"column:person;index;not null;size:32;"`
	Title       string                `gorm:"column:title;not null;size:200;"`
	Description string                `gorm:"column:description;type:text;"`
	Media       datatypes.JSON        `gorm:"column:media;type:jsonb;"`
}

func (m *FbM) TableName() string {
	return "hyper_feedback"
}
