package main

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyper/warehouse/pwh"
	"github.com/spf13/cast"
	"time"
)

func main() {
	helix.AfterStartup(func() {
		var curPwhM *pwh.PhysicalWhM
		err := pwh.CreatePwh(
			context.Background(),
			collar.Build("wlk", cast.ToString(time.Now().Unix())),
			"ChannelA-PWH",
			func(ctx context.Context, pwhM *pwh.PhysicalWhM) error {
				fmt.Println(hjson.MustToString(pwhM))
				curPwhM = pwhM
				return nil
			},
		)
		if err != nil {
			panic(err)
		}
		err = pwh.Into(context.Background(), pwh.IntoOutParas{
			PwhID:    curPwhM.ID,
			SkuID:    111,
			Quantity: 1000,
			Price:    100,
			Meta:     dict.NewDict().Set("hello", "world"),
		})
		if err != nil {
			panic(err)
		}
		err = pwh.Out(context.Background(), pwh.IntoOutParas{
			PwhID:    curPwhM.ID,
			SkuID:    111,
			Quantity: 100,
			Price:    100,
			Meta:     nil,
		})
		if err != nil {
			panic(err)
		}
		err = pwh.Lock(context.Background(), pwh.LockUnlockParas{
			PwhID:    curPwhM.ID,
			SkuID:    111,
			Quantity: 200,
			Meta:     nil,
		})
		if err != nil {
			panic(err)
		}

		err = pwh.Unlock(context.Background(), pwh.LockUnlockParas{
			PwhID:    curPwhM.ID,
			SkuID:    111,
			Quantity: 199,
			Meta:     nil,
		})
		if err != nil {
			panic(err)
		}

		m, err := pwh.GetSku(context.Background(), curPwhM.ID, 111)
		if err != nil {
			panic(err)
		}
		fmt.Println(hjson.MustToString(m))
	})
	helix.Startup()
}
