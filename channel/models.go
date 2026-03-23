package channel

import (
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hypes/collar"
	"gorm.io/datatypes"
)

type ChnM struct {
	hdb.Basic
	Biz       collar.ID      `gorm:"column:biz;index;size:64;uniqueIndex:uk_parent_name_biz;"`
	ID        ID             `gorm:"column:id;primaryKey;autoIncrement:false;"`
	Parent    ID             `gorm:"column:parent;uniqueIndex:uk_parent_name_biz;"`
	Name      string         `gorm:"column:name;uniqueIndex:uk_parent_name_biz;not null;size:32;"`
	Icon      string         `gorm:"column:icon;size:300;"`
	Seq       int            `gorm:"column:seq;"`
	Meta      datatypes.JSON `gorm:"column:meta;type:jsonb;"`
	Available bool           `gorm:"column:available"`
}

func (m *ChnM) TableName() string {
	return "hyper_channel"
}

func (m *ChnM) ToChannel() *Channel {
	meta := dict.NewDict()
	if len(m.Meta) > 0 {
		meta = *hjson.MustFromBytes[dict.Dict](m.Meta)
		if meta == nil {
			meta = dict.NewDict()
		}
	}
	return &Channel{
		ID:        m.ID,
		Name:      m.Name,
		Icon:      m.Icon,
		Seq:       m.Seq,
		Meta:      meta,
		Children:  make([]*Channel, 0),
		Available: m.Available,
	}
}
