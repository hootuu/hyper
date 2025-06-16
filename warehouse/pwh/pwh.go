package pwh

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/storage/hpg"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/collar"
	"go.uber.org/zap"
)

func CreatePwh(
	ctx context.Context,
	collar collar.Collar,
	memo string,
	call func(ctx context.Context, pwhM *PhysicalWhM) error,
) error {
	tx := hpg.CtxTx(ctx)
	if tx == nil {
		tx = zplt.HelixPgDB().PG()
	}
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

func Into(
	ctx context.Context,
	skuID uint64,
	quantity uint64,
	price uint64,
	memo string,
	link collar.Collar,
)
