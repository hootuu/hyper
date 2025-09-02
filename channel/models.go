package channel

import (
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hypes/collar"
)

type ChnM struct {
	hdb.Basic
	Biz       collar.ID `gorm:"column:biz;index;size:64;"`
	ID        ID        `gorm:"column:id;primaryKey;autoIncrement:false;"`
	Parent    ID        `gorm:"column:parent;uniqueIndex:uk_parent_name;"`
	Name      string    `gorm:"column:name;uniqueIndex:uk_parent_name;not null;size:32;"`
	Icon      string    `gorm:"column:icon;size:300;"`
	Seq       int       `gorm:"column:seq;"`
	Available bool      `gorm:"column:available"`
}

func (m *ChnM) TableName() string {
	return "hyper_channel"
}

func (m *ChnM) ToChannel() *Channel {
	return &Channel{
		ID:        m.ID,
		Name:      m.Name,
		Icon:      m.Icon,
		Seq:       m.Seq,
		Children:  make([]*Channel, 0),
		Available: m.Available,
	}
}
