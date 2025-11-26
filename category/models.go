package category

import (
	"github.com/hootuu/helix/storage/hdb"
)

type CtgM struct {
	hdb.Basic
	ID     ID     `gorm:"column:id;primaryKey;autoIncrement:false;"`
	Biz    string `gorm:"column:biz;not null;size:32;uniqueIndex:uk_parent_name;"`
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
