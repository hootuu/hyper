package category

import (
	"github.com/hootuu/helix/components/htree"
)

type ID = htree.ID

const Root ID = 0

type Categ struct {
	ID         ID       `json:"id"`
	Name       string   `json:"name"`
	Icon       string   `json:"icon"`
	CreateTime string   `json:"create_time"`
	Children   []*Categ `json:"children"`
}

func (c *Categ) AddChild(child *Categ) *Categ {
	c.Children = append(c.Children, child)
	return c
}
