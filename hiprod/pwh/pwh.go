package pwh

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyper/hiprod/prod"
	"github.com/hootuu/hyper/hyperplt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func CreatePwh(
	ctx context.Context,
	link collar.ID,
	memo string,
	ext dict.Dict,
) (ID, error) {
	tx := hyperplt.Tx(ctx)
	bExist, err := hdb.Exist[PhysicalWhM](tx, "link = ?", link)
	if err != nil {
		hlog.Err("hyper.pwh.CreatePwh: hdb.Exist[PhysicalWhM]", zap.Error(err))
		return 0, err
	}
	if bExist {
		return 0, fmt.Errorf("exist pwh: %s", link)
	}
	pwhM := &PhysicalWhM{
		Template: hdb.Template{
			Meta: hjson.MustToBytes(ext),
		},
		ID:   nextID(),
		Link: link,
		Memo: memo,
	}
	err = hdb.Create[PhysicalWhM](tx, pwhM)
	if err != nil {
		hlog.Err("hyper.pwh.CreatePwh: hdb.Create[PhysicalWhM]", zap.Error(err))
		return 0, err
	}
	return pwhM.ID, nil
}

type IntoOutParas struct {
	PwhID    ID         `json:"pwh_id"`
	SkuID    prod.SkuID `json:"sku_id"`
	Quantity uint64     `json:"quantity"`
	Price    uint64     `json:"price"`
	Meta     dict.Dict  `json:"meta"`
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

func Into(ctx context.Context, p IntoOutParas) (err error) {
	defer hlog.ElapseWithCtx(ctx, "pwh.Into", hlog.F(
		zap.Uint64("pwh", p.PwhID),
		zap.Uint64("sku", p.SkuID),
	),
		func() []zap.Field {
			if err != nil {
				return []zap.Field{zap.Error(err)}
			}
			return nil
		},
	)()
	if err := p.Validate(); err != nil {
		return err
	}
	tx := hyperplt.Tx(ctx)
	bPwhExist, err := hdb.Exist[PhysicalWhM](tx, "id = ?", p.PwhID)
	if err != nil {
		return err
	}
	if !bPwhExist {
		return fmt.Errorf("not exist pwh: %d", p.PwhID)
	}
	err = hdb.Tx(tx, func(tx *gorm.DB) error {
		pwhSkuM := &PhysicalSkuM{
			PWH:       p.PwhID,
			SKU:       p.SkuID,
			Available: 0,
			Locked:    0,
			Version:   0,
		}
		err = hdb.GetOrCreate[PhysicalSkuM](tx, pwhSkuM, "pwh = ? AND sku = ?", p.PwhID, p.SkuID)
		if err != nil {
			return err
		}
		mut := map[string]interface{}{
			"available": gorm.Expr("available + ?", p.Quantity),
			"version":   gorm.Expr("version + 1"),
		}
		err = hdb.Update[PhysicalSkuM](tx, mut,
			"pwh = ? AND sku = ? AND version = ?", p.PwhID, p.SkuID, pwhSkuM.Version)
		if err != nil {
			return err
		}
		err = hdb.Create[PhysicalInOutM](tx, &PhysicalInOutM{
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
	tx := hyperplt.Tx(ctx)
	bPwhExist, err := hdb.Exist[PhysicalWhM](tx, "id = ?", p.PwhID)
	if err != nil {
		return err
	}
	if !bPwhExist {
		return fmt.Errorf("not exist pwh: %d", p.PwhID)
	}
	pwhSkuM, err := hdb.MustGet[PhysicalSkuM](tx, "pwh = ? AND sku = ?", p.PwhID, p.SkuID)
	if err != nil {
		return err
	}
	if pwhSkuM.Available < p.Quantity {
		return fmt.Errorf("no available quantity")
	}
	err = hdb.Tx(tx, func(tx *gorm.DB) error {
		mut := map[string]interface{}{
			"available": gorm.Expr("available - ?", p.Quantity),
			"version":   gorm.Expr("version + 1"),
		}
		err = hdb.Update[PhysicalSkuM](tx, mut,
			"pwh = ? AND sku = ? AND version = ?", p.PwhID, p.SkuID, pwhSkuM.Version)
		if err != nil {
			return err
		}
		err = hdb.Create[PhysicalInOutM](tx, &PhysicalInOutM{
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

type LockUnlockParas struct {
	PwhID    ID        `json:"pwh_id"`
	SkuID    uint64    `json:"sku_id"`
	Quantity uint64    `json:"quantity"`
	Meta     dict.Dict `json:"meta"`
}

func (p LockUnlockParas) Validate() error {
	if p.PwhID == 0 {
		return fmt.Errorf("pwh_id is required")
	}
	if p.SkuID == 0 {
		return fmt.Errorf("require sku id")
	}
	if p.Quantity == 0 {
		return fmt.Errorf("require quantity")
	}
	return nil
}

func Lock(ctx context.Context, p LockUnlockParas) error {
	if err := p.Validate(); err != nil {
		return err
	}
	tx := hyperplt.Tx(ctx)
	bPwhExist, err := hdb.Exist[PhysicalWhM](tx, "id = ?", p.PwhID)
	if err != nil {
		return err
	}
	if !bPwhExist {
		return fmt.Errorf("not exist pwh: %d", p.PwhID)
	}
	pwhSkuM, err := hdb.MustGet[PhysicalSkuM](tx, "pwh = ? AND sku = ?", p.PwhID, p.SkuID)
	if err != nil {
		return err
	}
	if pwhSkuM.Available < p.Quantity {
		return fmt.Errorf("no available quantity")
	}
	err = hdb.Tx(tx, func(tx *gorm.DB) error {
		mut := map[string]interface{}{
			"available": gorm.Expr("available - ?", p.Quantity),
			"locked":    gorm.Expr("locked + ?", p.Quantity),
			"version":   gorm.Expr("version + 1"),
		}
		err = hdb.Update[PhysicalSkuM](tx, mut,
			"pwh = ? AND sku = ? AND version = ?", p.PwhID, p.SkuID, pwhSkuM.Version)
		if err != nil {
			return err
		}
		err = hdb.Create[PhysicalLockUnlockM](tx, &PhysicalLockUnlockM{
			PWH:       p.PwhID,
			SKU:       p.SkuID,
			Direction: DirectionLock,
			Quantity:  p.Quantity,
			Meta:      hjson.MustToBytes(p.Meta),
		})
		return nil
	})
	return nil
}

func Unlock(ctx context.Context, p LockUnlockParas) error {
	if err := p.Validate(); err != nil {
		return err
	}
	tx := hyperplt.Tx(ctx)
	bPwhExist, err := hdb.Exist[PhysicalWhM](tx, "id = ?", p.PwhID)
	if err != nil {
		return err
	}
	if !bPwhExist {
		return fmt.Errorf("not exist pwh: %d", p.PwhID)
	}
	pwhSkuM, err := hdb.MustGet[PhysicalSkuM](tx, "pwh = ? AND sku = ?", p.PwhID, p.SkuID)
	if err != nil {
		return err
	}
	if pwhSkuM.Locked < p.Quantity {
		return fmt.Errorf("no locked quantity")
	}
	err = hdb.Tx(tx, func(tx *gorm.DB) error {
		mut := map[string]interface{}{
			"available": gorm.Expr("available + ?", p.Quantity),
			"locked":    gorm.Expr("locked - ?", p.Quantity),
			"version":   gorm.Expr("version + 1"),
		}
		err = hdb.Update[PhysicalSkuM](tx, mut,
			"pwh = ? AND sku = ? AND version = ?", p.PwhID, p.SkuID, pwhSkuM.Version)
		if err != nil {
			return err
		}
		err = hdb.Create[PhysicalLockUnlockM](tx, &PhysicalLockUnlockM{
			PWH:       p.PwhID,
			SKU:       p.SkuID,
			Direction: DirectionUnlock,
			Quantity:  p.Quantity,
			Meta:      hjson.MustToBytes(p.Meta),
		})
		return nil
	})
	return nil
}

func GetSku(ctx context.Context, pwhID ID, skuID uint64) (*PhysicalSkuM, error) {
	tx := hyperplt.Tx(ctx)
	pwhSkuM, err := hdb.Get[PhysicalSkuM](tx, "pwh = ? AND sku = ?", pwhID, skuID)
	if err != nil {
		return nil, err
	}
	return pwhSkuM, nil
}
