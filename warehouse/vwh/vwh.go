package vwh

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyper/sku"
	"github.com/hootuu/hyper/warehouse/pwh"
	"github.com/hootuu/hyper/warehouse/vwh/strategy"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func CreateVwh(
	ctx context.Context,
	collar collar.Collar,
	meta dict.Dict,
	call func(ctx context.Context, vwhM *VirtualWhM) error,
) error {
	tx := db(ctx)
	bExist, err := hdb.Exist[VirtualWhM](tx, "collar = ?", collar)
	if err != nil {
		hlog.Err("hyper.vwh.CreateVwh: hdb.Exist[VirtualWhM]", zap.Error(err))
		return err
	}
	if bExist {
		return fmt.Errorf("exist pwh: %s", collar)
	}
	vwhM := &VirtualWhM{
		ID:     gVwhIdGenerator.NextUint64(),
		Collar: collar.String(),
		Meta:   hjson.MustToBytes(meta),
	}
	err = hdb.Create[VirtualWhM](tx, vwhM)
	if err != nil {
		hlog.Err("hyper.vwh.CreateVwh: hdb.Create[VirtualWhM]", zap.Error(err))
		return err
	}
	err = call(ctx, vwhM)
	if err != nil {
		return err
	}
	return nil
}

type AddPwhSrcParas struct {
	Vwh       ID                 `json:"vwh"`
	Pwh       pwh.ID             `json:"pwh"`
	Pricing   strategy.Pricing   `json:"price"`
	Inventory strategy.Inventory `json:"inventory"`
}

func (p AddPwhSrcParas) Validate() error {
	if p.Vwh == 0 {
		return fmt.Errorf("require vwh")
	}
	if p.Pwh == 0 {
		return fmt.Errorf("require pwh")
	}
	if err := p.Pricing.Validate(); err != nil {
		return err
	}
	if err := p.Inventory.Validate(); err != nil {
		return err
	}
	return nil
}

func AddPwhSrc(ctx context.Context, paras AddPwhSrcParas) error {
	if err := paras.Validate(); err != nil {
		return err
	}

	tx := db(ctx)

	if err := pwh.MustExist(ctx, paras.Pwh); err != nil {
		return err
	}
	hasThis, err := hdb.Exist[VirtualWhSrcM](tx,
		"vwh = ? AND pwh = ?", paras.Vwh, paras.Pwh)
	if err != nil {
		return err
	}
	if hasThis {
		return fmt.Errorf("exist vwh:%d, pwh: %d", paras.Vwh, paras.Pwh)
	}
	m := &VirtualWhSrcM{
		Vwh:           paras.Vwh,
		Pwh:           paras.Pwh,
		PricingType:   paras.Pricing.GetType(),
		InventoryType: paras.Inventory.GetType(),
		Pricing:       paras.Pricing.ToBytes(),
		Inventory:     paras.Inventory.ToBytes(),
	}
	err = hdb.Create[VirtualWhSrcM](tx, m)
	if err != nil {
		hlog.Err("hyper.vwh.AddPwhSrc: hdb.Create[VirtualWhSrcM]", zap.Error(err))
		return err
	}
	return nil
}

type SetSkuParas struct {
	Vwh       ID     `json:"vwh"`
	Sku       sku.ID `json:"sku"`
	Pwh       pwh.ID `json:"pwh"`
	Price     uint64 `json:"price"`
	Inventory uint64 `json:"inventory"`
}

func (p SetSkuParas) Validate() error {
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

func SetSku(ctx context.Context, paras SetSkuParas) error {
	if err := paras.Validate(); err != nil {
		return err
	}
	tx := db(ctx)
	err := hdb.GetOrCreate(tx, &VirtualWhSrcM{
		Vwh:           paras.Vwh,
		Pwh:           paras.Pwh,
		PricingType:   strategy.PricingProfit,
		InventoryType: strategy.InventoryTransfer,
		Pricing:       hjson.MustToBytes(strategy.DefaultPricing()),
		Inventory:     hjson.MustToBytes(strategy.DefaultInventory()),
	})
	if err != nil {
		return err
	}
	vwhSkuM, err := hdb.Get[VirtualWhSkuM](tx, "vwh = ? AND sku = ? AND pwh = ?",
		paras.Vwh, paras.Sku, paras.Pwh)
	if err != nil {
		return err
	}
	if vwhSkuM == nil {
		err = hdb.Create[VirtualWhSkuM](tx, &VirtualWhSkuM{
			Vwh:       paras.Vwh,
			Sku:       paras.Sku,
			Pwh:       paras.Pwh,
			Price:     paras.Price,
			Inventory: paras.Inventory,
			Version:   0,
		})
		if err != nil {
			return err
		}
		return nil
	}
	mut := make(map[string]any)
	if vwhSkuM.Price != paras.Price {
		mut["price"] = paras.Price
	}
	if vwhSkuM.Inventory != paras.Inventory {
		mut["inventory"] = paras.Inventory
	}
	if len(mut) == 0 {
		return nil
	}
	mut["version"] = gorm.Expr("version + 1")
	err = hdb.Update[VirtualWhSkuM](tx, mut, "vwh = ? AND sku = ? AND pwh = ?",
		paras.Vwh, paras.Sku, paras.Pwh)
	if err != nil {
		return err
	}
	return nil
}

func db(ctx context.Context) *gorm.DB {
	tx := hdb.CtxTx(ctx)
	if tx == nil {
		tx = zplt.HelixPgDB().PG()
	}
	return tx
}
