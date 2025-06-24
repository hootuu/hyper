package loader

import (
	"github.com/hootuu/hyper/sku"
	"github.com/hootuu/hyper/warehouse/pwh"
)

type Next byte

const (
	NextContinue Next = 1
	NextBreak    Next = 0
)

type SkuItem struct {
	Sku       sku.ID `json:"sku"`
	Pwh       pwh.ID `json:"pwh"`
	Price     uint64 `json:"price"`
	Inventory uint64 `json:"inventory"`
}

type Loader interface {
	Load(call func(item *SkuItem) Next) error
}
