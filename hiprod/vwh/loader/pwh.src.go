package loader

import "github.com/hootuu/hyper/warehouse/vwh/strategy"

type PwhSrcLoader struct {
	pricing   *strategy.Pricing
	inventory *strategy.Inventory
}

func NewPwhSrcLoader(pricing *strategy.Pricing, inventory *strategy.Inventory) *PwhSrcLoader {
	return &PwhSrcLoader{
		pricing:   pricing,
		inventory: inventory,
	}
}

func (loader *PwhSrcLoader) Load(call func(item *SkuItem) Next) error {
	return nil
}
