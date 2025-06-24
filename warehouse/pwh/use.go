package pwh

import (
	"context"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
)

func init() {
	helix.Use(helix.BuildHelix("hyper_warehouse_pwh", func() (context.Context, error) {
		err := zplt.HelixPgDB().PG().AutoMigrate(
			&PhysicalWhM{},
			&PhysicalSkuM{},
			&PhysicalInOutM{},
			&PhysicalLockUnlockM{},
		)
		if err != nil {
			return nil, err
		}
		err = initPwhIdGenerator()
		if err != nil {
			return nil, err
		}
		return nil, nil
	}, func(ctx context.Context) {

	}))
}
