package pwh

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/storage/hpg"
)

func MustExist(ctx context.Context, id ID) error {
	b, err := hpg.Exist[PhysicalWhM](db(ctx), id)
	if err != nil {
		return err
	}
	if !b {
		return fmt.Errorf("no such pwh: %d", id)
	}
	return nil
}
