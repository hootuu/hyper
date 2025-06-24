package channel

import (
	"context"
	"github.com/hootuu/helix/helix"
)

func init() {
	helix.Use(helix.BuildHelix("hyper_channel", func() (context.Context, error) {
		var err error
		err = db(nil).AutoMigrate(
			&ChnM{},
		)
		if err != nil {
			return nil, err
		}
		err = initChannelIdTree()
		if err != nil {
			return nil, err
		}
		return nil, nil
	}, func(ctx context.Context) {

	}))
}
