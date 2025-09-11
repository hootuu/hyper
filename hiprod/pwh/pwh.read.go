package pwh

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyper/hyperplt"
	"github.com/spf13/cast"
)

func BuildPwhCollar(pwhID ID) collar.Collar {
	return collar.Build("PWH", cast.ToString(pwhID))
}

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

func GetByLink(link collar.ID) ID {
	get, err := hdb.MustGet[PhysicalWhM](hyperplt.DB(), "link = ?", link)
	if err != nil {
		return 0
	}
	return get.ID
}

func MustGetById(Id ID) (*PhysicalWhM, error) {
	return hdb.MustGet[PhysicalWhM](hyperplt.DB(), "id = ?", Id)
}
