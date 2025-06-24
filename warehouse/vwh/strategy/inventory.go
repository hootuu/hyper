package strategy

import (
	"fmt"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hmath"
)

type Inventory interface {
	GetType() string
	Validate() error
	ToBytes() []byte
}

const InventoryTransfer = "transfer"

type TransferInventory struct {
	Rate hmath.Rate `json:"rate"`
}

func (t *TransferInventory) GetType() string {
	return InventoryTransfer
}

func (t *TransferInventory) Validate() error {
	if t.Rate.Rate() > 1 {
		return fmt.Errorf("rate too big")
	}
	return nil
}

func (t *TransferInventory) ToBytes() []byte {
	return hjson.MustToBytes(t)
}

const InventoryLock = "lock"

type LockInventory struct {
	Rate hmath.Rate `json:"rate"`
}

func (t *LockInventory) GetType() string {
	return InventoryLock
}

func (t *LockInventory) Validate() error {
	if t.Rate.Rate() > 1 {
		return fmt.Errorf("rate too big")
	}
	return nil
}

func (t *LockInventory) ToBytes() []byte {
	return hjson.MustToBytes(t)
}

func DefaultInventory() Inventory {
	return &TransferInventory{hmath.NewRate(hmath.Rate10000, 10000)}
}
