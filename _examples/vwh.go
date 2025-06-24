package main

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyper/warehouse/pwh"
	"github.com/hootuu/hyper/warehouse/vwh"
	"github.com/hootuu/hyper/warehouse/vwh/strategy"
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
		var curVwhM *vwh.VirtualWhM
		err = vwh.CreateVwh(context.Background(),
			collar.Build("wlk", cast.ToString(time.Now().Unix())),
			dict.NewDict().Set("channel", "ABC"),
			func(ctx context.Context, vwhM *vwh.VirtualWhM) error {
				fmt.Println(hjson.MustToString(vwhM))
				curVwhM = vwhM
				return nil
			})
		if err != nil {
			panic(err)
		}
		err = vwh.AddPwhSrc(context.Background(), vwh.AddPwhSrcParas{
			Vwh:       curVwhM.ID,
			Pwh:       curPwhM.ID,
			Pricing:   strategy.DefaultPricing(),
			Inventory: strategy.DefaultInventory(),
		})
		if err != nil {
			panic(err)
		}
		err = vwh.SetSku(context.Background(), vwh.SetSkuParas{
			Vwh:       curVwhM.ID,
			Sku:       123456,
			Pwh:       curPwhM.ID,
			Price:     600,
			Inventory: 100,
		})
		if err != nil {
			panic(err)
		}
		err = vwh.SetSku(context.Background(), vwh.SetSkuParas{
			Vwh:       curVwhM.ID,
			Sku:       123456,
			Pwh:       curPwhM.ID,
			Price:     700,
			Inventory: 111,
		})
		if err != nil {
			panic(err)
		}
	})
	helix.Startup()
}
