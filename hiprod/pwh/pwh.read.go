package pwh

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyper/hyperplt"
)

func ExistByLink(ctx context.Context, link collar.ID) (bool, error) {
	tx := hyperplt.Tx(ctx)
	return hdb.Exist[PhysicalWhM](tx, "link = ?", link)
}

func Exist(ctx context.Context, id ID) (bool, error) {
	tx := hyperplt.Tx(ctx)
	return hdb.Exist[PhysicalWhM](tx, "id = ?", id)
}

func MustExist(ctx context.Context, id ID) error {
	b, err := Exist(ctx, id)
	if err != nil {
		return err
	}
	if !b {
		return fmt.Errorf("no such pwh: %d", id)
	}
	return nil
}
