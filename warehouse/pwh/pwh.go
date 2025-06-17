package pwh

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/storage/hpg"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/collar"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func CreatePwh(
	ctx context.Context,
	collar collar.Collar,
	memo string,
	call func(ctx context.Context, pwhM *PhysicalWhM) error,
) error {
	tx := db(ctx)
	bExist, err := hpg.Exist[PhysicalWhM](tx, "collar = ?", collar)
	if err != nil {
		hlog.Err("hyper.pwh.CreatePwh: hpg.Exist[PhysicalWhM]", zap.Error(err))
		return err
	}
	if bExist {
		return fmt.Errorf("exist pwh: %s", collar)
	}
	pwhM := &PhysicalWhM{
		ID:     gPwhIdGenerator.NextUint64(),
		Collar: collar.String(),
		Memo:   memo,
	}
	err = hpg.Create[PhysicalWhM](tx, pwhM)
	if err != nil {
		hlog.Err("hyper.pwh.CreatePwh: hpg.Create[PhysicalWhM]", zap.Error(err))
		return err
	}
	err = call(ctx, pwhM)
	if err != nil {
		return err
	}
	return nil
}

type IntoOutParas struct {
	PwhID    ID        `json:"pwh_id"`
	SkuID    uint64    `json:"sku_id"`
	Quantity uint64    `json:"quantity"`
	Price    uint64    `json:"price"`
	Meta     dict.Dict `json:"meta"`
}

func (p IntoOutParas) Validate() error {
	if p.PwhID == 0 {
		return fmt.Errorf("pwh_id is required")
	}
	if p.SkuID == 0 {
		return fmt.Errorf("require sku id")
	}
	if p.Quantity == 0 {
		return fmt.Errorf("require quantity")
	}
	if p.Price == 0 {
		return fmt.Errorf("require price")
	}
	return nil
}

func Into(ctx context.Context, p IntoOutParas) error {
	if err := p.Validate(); err != nil {
		return err
	}
	tx := db(ctx)
	bPwhExist, err := hpg.Exist[PhysicalWhM](tx, "id = ?", p.PwhID)
	if err != nil {
		return err
	}
	if !bPwhExist {
		return fmt.Errorf("not exist pwh: %d", p.PwhID)
	}
	err = hpg.Tx(tx, func(tx *gorm.DB) error {
		pwhSkuM := &PhysicalSkuM{
			PWH:       p.PwhID,
			SKU:       p.SkuID,
			Available: 0,
			Locked:    0,
			Version:   0,
		}
		err = hpg.GetOrCreate[PhysicalSkuM](tx, pwhSkuM, "pwh = ? AND sku = ?", p.PwhID, p.SkuID)
		if err != nil {
			return err
		}
		mut := map[string]interface{}{
			"available": gorm.Expr("available + ?", p.Quantity),
			"version":   gorm.Expr("version + 1"),
		}
		err = hpg.Update[PhysicalSkuM](tx, mut,
			"pwh = ? AND sku = ? AND version = ?", p.PwhID, p.SkuID, pwhSkuM.Version)
		if err != nil {
			return err
		}
		err = hpg.Create[PhysicalInOutM](tx, &PhysicalInOutM{
			PWH:       p.PwhID,
			SKU:       p.SkuID,
			Direction: DirectionIn,
			Quantity:  p.Quantity,
			Price:     p.Price,
			Meta:      hjson.MustToBytes(p.Meta),
		})
		return nil
	})
	return nil
}

func Out(ctx context.Context, p IntoOutParas) error {
	if err := p.Validate(); err != nil {
		return err
	}
	tx := db(ctx)
	bPwhExist, err := hpg.Exist[PhysicalWhM](tx, "id = ?", p.PwhID)
	if err != nil {
		return err
	}
	if !bPwhExist {
		return fmt.Errorf("not exist pwh: %d", p.PwhID)
	}
	pwhSkuM, err := hpg.MustGet[PhysicalSkuM](tx, "pwh = ? AND sku = ?", p.PwhID, p.SkuID)
	if err != nil {
		return err
	}
	if pwhSkuM.Available < p.Quantity {
		return fmt.Errorf("no available quantity")
	}
	err = hpg.Tx(tx, func(tx *gorm.DB) error {
		mut := map[string]interface{}{
			"available": gorm.Expr("available - ?", p.Quantity),
			"version":   gorm.Expr("version + 1"),
		}
		err = hpg.Update[PhysicalSkuM](tx, mut,
			"pwh = ? AND sku = ? AND version = ?", p.PwhID, p.SkuID, pwhSkuM.Version)
		if err != nil {
			return err
		}
		err = hpg.Create[PhysicalInOutM](tx, &PhysicalInOutM{
			PWH:       p.PwhID,
			SKU:       p.SkuID,
			Direction: DirectionOut,
			Quantity:  p.Quantity,
			Price:     p.Price,
			Meta:      hjson.MustToBytes(p.Meta),
		})
		return nil
	})
	return nil
}

func db(ctx context.Context) *gorm.DB {
	tx := hpg.CtxTx(ctx)
	if tx == nil {
		tx = zplt.HelixPgDB().PG()
	}
	return tx
}
