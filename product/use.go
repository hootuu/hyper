package product

import (
	"context"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
)

func init() {
	helix.Use(helix.BuildHelix("hyper_product", func() (context.Context, error) {
		err := zplt.HelixPgDB().PG().AutoMigrate(
			&SpuM{},
		)
		if err != nil {
			return nil, err
		}
		err = initSpuIdGenerator()
		if err != nil {
			return nil, err
		}
		return nil, nil
	}, func(ctx context.Context) {

	}))
}
