package vwh

import (
	"context"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
)

func init() {
	helix.Use(helix.BuildHelix("hyper_warehouse_vwh", func() (context.Context, error) {
		err := zplt.HelixPgDB().PG().AutoMigrate(
			&VirtualWhM{},
			&VirtualWhSrcM{},
			&VirtualWhSkuM{},
		)
		if err != nil {
			return nil, err
		}
		err = initVwhIdGenerator()
		if err != nil {
			return nil, err
		}
		if err = initUni(); err != nil {
			return nil, err
		}
		return nil, nil
	}, func(ctx context.Context) {

	}))
}
