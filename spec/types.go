package spec

import (
	"github.com/hootuu/hyle/hypes/media"
	"github.com/hootuu/hyper/category"
)

type ID = int64
type OptID = int64

type Spec struct {
	ID       ID          `json:"id"`
	Category category.ID `json:"category"`
	Name     string      `json:"name"`
	Intro    string      `json:"intro,omitempty"`
}

type Option struct {
	OptID OptID      `json:"opt_id"`
	Label string     `json:"label"`
	Media media.More `json:"media,omitempty"`
	Seq   int        `json:"seq"`
}
