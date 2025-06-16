package product

import (
	"github.com/hootuu/hyle/hypes/media"
	"github.com/hootuu/hyper/spec"
)

type SpuID = string
type SkuID = string

type SpuSpecSetting struct {
	Spu   SpuID      `json:"spu"`
	Specs []*SpuSpec `json:"specs,omitempty"`
}

type SpuSpec struct {
	Spec    spec.ID        `json:"spec"`
	Options []*spec.Option `json:"options,omitempty"`
	Seq     int            `json:"seq"`
}

type SpecOpt struct {
	ID    spec.OptID `json:"id"`
	Label string     `json:"label"`
	Media media.More `json:"media,omitempty"`
	Seq   int        `json:"seq"`
}

type SkuSpecSetting struct {
	Spu   SpuID        `json:"spu"`
	Specs []spec.OptID `json:"specs,omitempty"`
}
