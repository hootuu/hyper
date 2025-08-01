package vwh

import (
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyper/hyperplt"
)

func DbVwhGet(id ID) (*VirtualWhM, error) {
	return hdb.Get[VirtualWhM](hyperplt.DB(), "id = ?", id)
}
