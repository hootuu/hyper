package prod

import (
	"github.com/hootuu/hyle/hypes/media"
)

const (
	UniBiz Biz = "UNI"
)

type Biz = string
type SpuID = uint64
type SkuID = uint64
type SpecID = uint64
type SpecOptID = uint64

type SpuSpecSetting struct {
	Spu   SpuID      `json:"spu"`
	Specs []*SpuSpec `json:"specs,omitempty"`
}

type SpuSpec struct {
	Spec    SpecID     `json:"spec"`
	Options []*SpecOpt `json:"options,omitempty"`
	Seq     int        `json:"seq"`
}

type SpecOpt struct {
	ID    SpecOptID  `json:"id"`
	Label string     `json:"label"`
	Media media.More `json:"media,omitempty"`
	Seq   int        `json:"seq"`
}

type SkuSpecSetting struct {
	Spu   SpuID       `json:"spu"`
	Specs []SpecOptID `json:"specs,omitempty"`
}
