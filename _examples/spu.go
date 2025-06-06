package main

import (
	"fmt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyle/hypes/media"
	"github.com/hootuu/hyper/product"
)

func main() {
	helix.AfterStartup(func() {
		spuM, err := product.CreateSpu(&product.SpuM{
			Collar:    collar.Build("shop", "123456"),
			Frontend:  0,
			Backend:   0,
			Name:      "SPU-TEST",
			Intro:     "SPU-INTRO",
			Brand:     "SPU-BRAND",
			Version:   0,
			MainMedia: hjson.MustToBytes(media.More{media.New(media.ImageType, "https://1.1")}),
			MoreMedia: nil,
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(spuM.ID)
	})
	helix.Startup()
}
