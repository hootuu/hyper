package spec

import (
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyper/category"
)

type CatSpecM struct {
	hdb.Basic
	ID       ID          `gorm:"column:id;primaryKey;"`
	Category category.ID `gorm:"column:category;uniqueIndex:uk_cat_name;"`
	Name     string      `gorm:"column:name;uniqueIndex:uk_cat_name;not null;size:32;"`
	Intro    string      `gorm:"column:intro;size:300;"`
}

func (m *CatSpecM) TableName() string {
	return "hyper_category_spec"
}

func (m *CatSpecM) To() *Spec {
	return &Spec{
		ID:       m.ID,
		Category: m.Category,
		Name:     m.Name,
		Intro:    m.Intro,
	}
}
