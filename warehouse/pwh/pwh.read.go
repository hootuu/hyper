package pwh

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/storage/hdb"
)

func MustExist(ctx context.Context, id ID) error {
	b, err := hdb.Exist[PhysicalWhM](db(ctx), id)
	if err != nil {
		return err
	}
	if !b {
		return fmt.Errorf("no such pwh: %d", id)
	}
	return nil
}
