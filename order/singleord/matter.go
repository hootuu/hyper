package singleord

import (
	"fmt"
	"github.com/hootuu/hyper/hiprod/prod"
	"github.com/hootuu/hyper/hiprod/vwh"
)

type Matter struct {
	SkuID    prod.SkuID `json:"sku_id"`
	VwhID    vwh.ID     `json:"vwh_id"`
	PwhID    vwh.ID     `json:"pwh_id"`
	Price    uint64     `json:"price"`
	Quantity uint32     `json:"quantity"`
	Amount   uint64     `json:"amount"`
}

func (m *Matter) Validate() error {
	if m.SkuID == 0 {
		return fmt.Errorf("require Matter.SkuID")
	}
	if m.Quantity == 0 {
		return fmt.Errorf("require Matter.Quantity")
	}
	return nil
}
