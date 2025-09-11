package channel

import (
	"fmt"
	"github.com/hootuu/helix/components/htree"
	"github.com/hootuu/hyle/hypes/collar"
)

type ID = htree.ID

const (
	CollarCode = "hyper_channel"
)

func Collar(id ID) collar.Collar {
	return collar.Build(CollarCode, fmt.Sprintf("%d", id))
}

const Root ID = 0

type Channel struct {
	ID        ID         `json:"id"`
	Name      string     `json:"name"`
	Icon      string     `json:"icon"`
	Seq       int        `json:"seq"`
	Children  []*Channel `json:"children"`
	Available bool       `gorm:"column:available"`
}

func (c *Channel) AddChild(child *Channel) *Channel {
	c.Children = append(c.Children, child)
	return c
}
