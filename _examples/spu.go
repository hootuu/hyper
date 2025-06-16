package main

import (
	"fmt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyle/hypes/media"
	"github.com/hootuu/hyper/product"
	"github.com/hootuu/hyper/spec"
	"github.com/spf13/cast"
	"time"
)

func main() {
	helix.AfterStartup(func() {
		defer hlog.Elapse("productTest")()
		for i := 0; i < 1000; i++ {
			productTest()
		}
	})
	helix.Startup()
}

func productTest() {
	spuM, err := product.CreateSpu(&product.SpuM{
		Collar:    collar.Build("shop", "123456"),
		Category:  0,
		Name:      "SPU-TEST-" + cast.ToString(time.Now().Unix()),
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
	s, err := product.CreateSpec(&product.SpuSpecSetting{
		Spu: spuM.ID,
		Specs: []*product.SpuSpec{
			{
				Spec: 999,
				Options: []*spec.Option{
					{
						//OptID: 0,
						Label: "RED SALES",
						Media: nil,
					},
					{
						//OptID: 0,
						Label: "BLUE SALES",
						Media: nil,
					},
				},
				//Seq: 0,
			},
			{
				Spec: 998,
				Options: []*spec.Option{
					{
						//OptID: 0,
						Label: "100 SALES",
						Media: nil,
					},
					{
						//OptID: 0,
						Label: "500 SALES",
						Media: nil,
					},
				},
				//Seq: 0,
			},
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(hjson.MustToString(s))
	var specs []spec.OptID
	for _, spuSpecItem := range s.Specs {
		for _, optItem := range spuSpecItem.Options {
			specs = append(specs, optItem.OptID)
		}
	}
	skuSetting := &product.SkuSpecSetting{
		Spu:   spuM.ID,
		Specs: specs,
	}
	skuID, err := product.CreateSku(skuSetting)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(hjson.MustToString(skuSetting))
	fmt.Println(skuID)
}
