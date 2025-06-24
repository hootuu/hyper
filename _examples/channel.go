package main

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyper/channel"
	"github.com/spf13/cast"
	"time"
)

func main() {
	helix.AfterStartup(func() {
		var rootM *channel.ChnM
		err := channel.Create(
			context.Background(),
			collar.Build("shop", cast.ToString(time.Now().Unix())),
			"channel"+cast.ToString(time.Now().Unix()),
			"",
			1,
			func(ctx context.Context, chnM *channel.ChnM) error {
				fmt.Println(hjson.MustToString(chnM))
				rootM = chnM
				return nil
			},
		)
		if err != nil {
			panic(err)
		}
		var level2M *channel.ChnM
		err = channel.Add(context.Background(), rootM.ID,
			"channel-1"+cast.ToString(time.Now().Unix()), "", 1,
			func(ctx context.Context, chnM *channel.ChnM) error {
				fmt.Println(hjson.MustToString(chnM))
				level2M = chnM
				return nil
			})
		if err != nil {
			panic(err)
		}

		var level3M *channel.ChnM
		err = channel.Add(context.Background(), level2M.ID,
			"channel-2"+cast.ToString(time.Now().Unix()), "", 1,
			func(ctx context.Context, chnM *channel.ChnM) error {
				fmt.Println(hjson.MustToString(chnM))
				level3M = chnM
				return nil
			})
		if err != nil {
			panic(err)
		}
		var level4M *channel.ChnM
		err = channel.Add(context.Background(), level3M.ID,
			"channel-3"+cast.ToString(time.Now().Unix()), "", 1,
			func(ctx context.Context, chnM *channel.ChnM) error {
				level4M = chnM
				fmt.Println(hjson.MustToString(level4M))
				return nil
			})
		if err != nil {
			panic(err)
		}
	})
	helix.Startup()
}
