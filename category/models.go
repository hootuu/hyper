package category

import (
	"github.com/hootuu/helix/storage/hpg"
)

type CtgM struct {
	hpg.Basic
	ID     ID     `gorm:"column:id;primaryKey;"`
	Parent ID     `gorm:"column:parent;uniqueIndex:uk_parent_name;"`
	Name   string `gorm:"column:name;uniqueIndex:uk_parent_name;not null;size:32;"`
	Icon   string `gorm:"column:icon;size:300;"`
}

func (m *CtgM) ToCateg() *Categ {
	return &Categ{
		ID:       m.ID,
		Name:     m.Name,
		Icon:     m.Icon,
		Children: make([]*Categ, 0),
	}
}
