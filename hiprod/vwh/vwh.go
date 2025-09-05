package vwh

import (
	"context"
	"fmt"

	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyper/hiprod/prod"
	"github.com/hootuu/hyper/hiprod/pwh"
	"github.com/hootuu/hyper/hiprod/vwh/strategy"
	"github.com/hootuu/hyper/hyperplt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func BuildVwhExtLink(vwh ID, sku prod.SkuID) string {
	return collar.Build("vwh_ext", fmt.Sprintf("%d_%d", vwh, sku)).ToID()
}

func CreateVwh(
	ctx context.Context,
	link collar.ID,
	memo string,
	ext dict.Dict,
) (ID, error) {
	tx := hyperplt.Tx(ctx)
	bExist, err := hdb.Exist[VirtualWhM](tx, "link = ?", link)
	if err != nil {
		hlog.Err("hyper.vwh.CreateVwh: hdb.Exist[VirtualWhM]", zap.Error(err))
		return 0, err
	}
	if bExist {
		return 0, fmt.Errorf("exist pwh: %s", link)
	}
	vwhM := &VirtualWhM{
		Template: hdb.Template{
			Meta: hjson.MustToBytes(ext),
		},
		ID:   gVwhIdGenerator.NextUint64(),
		Link: link,
		Memo: memo,
	}
	err = hdb.Create[VirtualWhM](tx, vwhM)
	if err != nil {
		hlog.Err("hyper.vwh.CreateVwh: hdb.Create[VirtualWhM]", zap.Error(err))
		return 0, err
	}
	return vwhM.ID, nil
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

	tx := hyperplt.Tx(ctx)

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
	Vwh       ID         `json:"vwh"`
	Sku       prod.SkuID `json:"sku"`
	Pwh       pwh.ID     `json:"pwh"`
	Price     uint64     `json:"price"`
	Inventory uint64     `json:"inventory"`
	Channel   uint64     `json:"channel"`
	UseExt    bool       `json:"use_ext"`
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
	tx := hyperplt.Tx(ctx)
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
		err = hdb.Tx(tx, func(innerTx *gorm.DB) error {
			err = hdb.Create[VirtualWhSkuM](innerTx, &VirtualWhSkuM{
				Vwh:          paras.Vwh,
				Sku:          paras.Sku,
				Pwh:          paras.Pwh,
				Price:        paras.Price,
				Inventory:    paras.Inventory,
				UseInventory: paras.Inventory > 0,
				Version:      0,
			})
			if err != nil {
				return err
			}
			if paras.UseExt {
				err = hdb.Create[VirtualWhSkuExtM](innerTx, &VirtualWhSkuExtM{
					Template:  hdb.Template{},
					Vwh:       paras.Vwh,
					Sku:       paras.Sku,
					Pwh:       paras.Pwh,
					Link:      BuildVwhExtLink(paras.Vwh, paras.Sku),
					Channel:   paras.Channel,
					Available: true,
					Sort:      0,
				})
				if err != nil {
					return err
				}
			}
			return nil
		})
		return err
	}
	_ = hdb.Update[VirtualWhSkuExtM](tx, map[string]any{
		"available": true,
	}, "vwh = ? AND sku = ? AND pwh = ?", paras.Vwh, paras.Sku, paras.Pwh)

	mut := make(map[string]any)
	if vwhSkuM.Price != paras.Price {
		mut["price"] = paras.Price
	}
	if paras.Inventory > 0 {
		mut["inventory"] = gorm.Expr("inventory + ?", paras.Inventory)
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

func DeductInventory(ctx context.Context, vwhID ID, skuID prod.SkuID, pwhID pwh.ID, quantity uint64) error {
	tx := hyperplt.Tx(ctx)
	vwhSkuM, err := hdb.MustGet[VirtualWhSkuM](tx, "vwh = ? AND sku = ? AND pwh = ?", vwhID, skuID, pwhID)
	if err != nil {
		return err
	}
	if vwhSkuM.Inventory < quantity {
		return fmt.Errorf("not enough inventory: %d < %d", vwhSkuM.Inventory, quantity)
	}
	mut := map[string]any{
		"inventory": gorm.Expr("inventory - ?", quantity),
		"version":   gorm.Expr("version + 1"),
	}
	err = hdb.Update[VirtualWhSkuM](tx, mut, "vwh = ? AND sku = ? AND pwh = ?",
		vwhID, skuID, pwhID)
	if err != nil {
		return err
	}
	return nil
}

func AddInventory(ctx context.Context, vwhID ID, skuID prod.SkuID, pwhID pwh.ID, quantity uint64) error {
	tx := hyperplt.Tx(ctx)
	vwhSkuM, err := hdb.MustGet[VirtualWhSkuM](tx, "vwh = ? AND sku = ? AND pwh = ?", vwhID, skuID, pwhID)
	if err != nil {
		return err
	}

	if vwhSkuM == nil {
		return fmt.Errorf("vwh sku not found: vwh=%d, sku=%d, pwh=%d", vwhID, skuID, pwhID)
	}

	mut := map[string]any{
		"inventory": gorm.Expr("inventory + ?", quantity),
		"version":   gorm.Expr("version + 1"),
	}
	err = hdb.Update[VirtualWhSkuM](tx, mut, "vwh = ? AND sku = ? AND pwh = ?",
		vwhID, skuID, pwhID)
	if err != nil {
		return err
	}
	return nil
}

type UpdateVwhSkuParas struct {
	Vwh     ID         `json:"vwh"`
	Sku     prod.SkuID `json:"sku"`
	Pwh     pwh.ID     `json:"pwh"`
	Price   uint64     `json:"price"`
	Channel uint64     `json:"channel"`
	Sort    uint64     `json:"sort"`
	Meta    dict.Dict  `json:"meta"`
}

func (p UpdateVwhSkuParas) Validate() error {
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

func UpdateSkuPrice(ctx context.Context, paras UpdateVwhSkuParas) error {
	if err := paras.Validate(); err != nil {
		return err
	}
	if paras.Price == 0 {
		return fmt.Errorf("require price")
	}
	tx := hyperplt.Tx(ctx)
	vwhSkuM, err := hdb.MustGet[VirtualWhSkuM](tx, "pwh = ? AND sku = ? And vwh = ?", paras.Pwh, paras.Sku, paras.Vwh)
	if err != nil {
		return err
	}
	if vwhSkuM.Price == paras.Price {
		return nil
	}
	mut := map[string]any{
		"price": paras.Price,
	}
	err = hdb.Update[VirtualWhSkuM](tx, mut, "auto_id = ?", vwhSkuM.AutoID)
	if err != nil {
		return err
	}
	return nil
}

func UpdateSkuExt(ctx context.Context, paras UpdateVwhSkuParas) error {
	if err := paras.Validate(); err != nil {
		return err
	}
	tx := hyperplt.Tx(ctx)
	vwhSkuExtM, err := hdb.MustGet[VirtualWhSkuExtM](tx, "pwh = ? AND link = ?", paras.Pwh, BuildVwhExtLink(paras.Vwh, paras.Sku))
	if err != nil {
		return err
	}
	mut := make(map[string]any)
	if paras.Sort > 0 {
		mut["sort"] = paras.Sort
	}
	if paras.Channel > 0 {
		mut["channel"] = paras.Channel
	}
	meta := *hjson.MustFromBytes[dict.Dict](vwhSkuExtM.Meta)
	if len(paras.Meta) > 0 {
		if len(meta) == 0 {
			meta = paras.Meta
		} else {
			for k, v := range paras.Meta {
				meta[k] = v
			}
		}
		mut["meta"] = hjson.MustToBytes(meta)
	}
	if len(mut) == 0 {
		return nil
	}
	err = hdb.Update[VirtualWhSkuExtM](tx, mut, "auto_id = ?", vwhSkuExtM.AutoID)
	if err != nil {
		return err
	}
	return nil
}
