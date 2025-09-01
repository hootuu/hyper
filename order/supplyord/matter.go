package supplyord

import (
	"fmt"
	"github.com/hootuu/hyle/hypes/ex"
	"github.com/hootuu/hyper/hiprod/prod"
	"github.com/hootuu/hyper/hiprod/vwh"
)

type Item struct {
	ProdID   prod.ID    `json:"id"`
	SkuID    prod.SkuID `json:"sku_id"`
	VwhID    vwh.ID     `json:"vwh_id"`
	PwhID    vwh.ID     `json:"pwh_id"`
	Price    uint64     `json:"price"`
	Quantity uint64     `json:"quantity"`
	Amount   uint64     `json:"amount"`
}

func (m *Item) Validate() error {
	if m.SkuID == 0 {
		return fmt.Errorf("require Matter.SkuID")
	}
	if m.Quantity == 0 {
		return fmt.Errorf("require Matter.Quantity")
	}
	return nil
}

type Matter struct {
	Items  []*Item `json:"items"`
	Amount uint64  `json:"amount"`
}

func (m Matter) GetDigest() ex.Meta {
	return ex.Meta{
		"items": m.Items,
	}
}

func (m Matter) Validate() error {
	if len(m.Items) == 0 {
		return fmt.Errorf("require Matter.Items")
	}
	if m.Amount == 0 {
		return fmt.Errorf("require Matter.Amount")
	}
	return nil
}
