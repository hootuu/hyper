package hiprod

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyper/hiprod/prod"
	"github.com/hootuu/hyper/hiprod/pwh"
	"github.com/hootuu/hyper/hiprod/vwh"
	"github.com/hootuu/hyper/hyperplt"
	"gorm.io/gorm"
	"time"
)

type SkuUnpPublishParas struct {
	Vwh  vwh.ID     `json:"vwh"`
	Sku  prod.SkuID `json:"sku"`
	Pwh  pwh.ID     `json:"pwh"`
	Meta dict.Dict  `json:"meta"`

	Quantity uint64 `json:"quantity"`
}

func (p SkuUnpPublishParas) Validate() error {
	if p.Vwh == 0 {
		return fmt.Errorf("require vwh")
	}
	if p.Pwh == 0 {
		return fmt.Errorf("require pwh")
	}
	if p.Sku == 0 {
		return fmt.Errorf("require sku")
	}
	return nil
}

func SkuUnpPublish(ctx context.Context, paras SkuUnpPublishParas) error {
	err := paras.Validate()
	if err != nil {
		return err
	}
	tx := hyperplt.Tx(ctx)
	vwhSkuM, err := hdb.Get[vwh.VirtualWhSkuM](tx, "vwh = ? AND sku = ? AND pwh = ?",
		paras.Vwh, paras.Sku, paras.Pwh)
	if err != nil {
		return err
	}
	if vwhSkuM == nil {
		return fmt.Errorf("vwh sku not found: vwh=%d, sku=%d, pwh=%d", paras.Vwh, paras.Sku, paras.Pwh)
	}
	vwhSkuExtM, err := hdb.Get[vwh.VirtualWhSkuExtM](tx, "pwh = ? AND link = ?", paras.Pwh, vwh.BuildVwhExtLink(paras.Vwh, paras.Sku))
	if err != nil {
		return err
	}

	err = hdb.Tx(tx, func(innerTx *gorm.DB) error {
		var inventory uint64
		quantity := vwhSkuM.Inventory
		if quantity > 0 {
			var useStock uint64
			if paras.Quantity > 0 {
				if paras.Quantity > quantity {
					useStock = paras.Quantity - quantity
				}
				if paras.Quantity < quantity {
					inventory = quantity - paras.Quantity
				}
				quantity = paras.Quantity
			}
			err = pwh.Unlock(hdb.TxCtx(innerTx), pwh.LockUnlockParas{
				PwhID:    paras.Pwh,
				SkuID:    paras.Sku,
				Quantity: quantity,
				Meta:     paras.Meta,
			})
			if err != nil {
				return err
			}
			if useStock > 0 {
				err = pwh.Out(hdb.TxCtx(innerTx), pwh.IntoOutParas{
					PwhID:    paras.Pwh,
					SkuID:    paras.Sku,
					Quantity: quantity,
					Price:    1,
					Meta:     paras.Meta,
				})
				if err != nil {
					return err
				}
			}
			err = hdb.Update[vwh.VirtualWhSkuM](innerTx, map[string]any{
				"inventory": inventory,
				"version":   gorm.Expr("version + 1"),
			}, "auto_id = ?", vwhSkuM.AutoID)
			if err != nil {
				return err
			}
		}
		if vwhSkuExtM != nil && inventory == 0 {
			err = hdb.Update[vwh.VirtualWhSkuExtM](innerTx, map[string]any{
				"available": false,
			}, "auto_id = ?", vwhSkuExtM.AutoID)
			if err != nil {
				return err
			}

			_ = hdb.Update[vwh.VirtualWhSkuM](innerTx, map[string]any{
				"updated_at": time.Now(),
			}, "auto_id = ?", vwhSkuM.AutoID)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

type SkuStockResetParas struct {
	Vwh      vwh.ID     `json:"vwh"`
	Sku      prod.SkuID `json:"sku"`
	Pwh      pwh.ID     `json:"pwh"`
	Quantity uint64     `json:"quantity"`
}

func (p SkuStockResetParas) Validate() error {
	if p.Vwh == 0 {
		return fmt.Errorf("require vwh")
	}
	if p.Pwh == 0 {
		return fmt.Errorf("require pwh")
	}
	if p.Sku == 0 {
		return fmt.Errorf("require sku")
	}
	if p.Quantity == 0 {
		return fmt.Errorf("require quantity")
	}
	return nil
}

func SkuStockReset(ctx context.Context, paras SkuStockResetParas) error {
	err := paras.Validate()
	if err != nil {
		return err
	}
	tx := hyperplt.Tx(ctx)
	vwhSkuM, err := hdb.MustGet[vwh.VirtualWhSkuM](tx, "vwh = ? AND sku = ? AND pwh = ?", paras.Vwh, paras.Sku, paras.Pwh)
	if err != nil {
		return err
	}
	if vwhSkuM.UseInventory {
		err = vwh.DeductInventory(ctx, paras.Vwh, paras.Sku, paras.Pwh, paras.Quantity)
	} else {
		err = pwh.Into(ctx, pwh.IntoOutParas{
			PwhID:    paras.Pwh,
			SkuID:    paras.Sku,
			Quantity: paras.Quantity,
			Price:    vwhSkuM.Price,
			Meta: map[string]interface{}{
				"into_out_biz": "ORD_RESET",
			},
		})
	}
	if err != nil {
		return err
	}
	return nil
}
