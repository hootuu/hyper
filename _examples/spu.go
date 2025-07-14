package main

import (
	"fmt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/data/pagination"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyle/hypes/media"
	"github.com/hootuu/hyper/hyperidx"
	"github.com/hootuu/hyper/product"
	"github.com/spf13/cast"
	"time"
)

func main() {
	helix.AfterStartup(func() {
		defer hlog.Elapse("productTest")()
		total := 1
		for i := 0; i < total; i++ {
			productTest()
		}

		s := time.Now().UnixMilli()
		for {
			time.Sleep(100 * time.Millisecond)
			filter := "category = 119"
			data, err := hyperidx.Filter(hyperidx.SpuIndex,
				filter, []string{}, pagination.NewPage(1, 1))
			if err != nil {
				return
			}
			fmt.Println(data.Paging)
			//fmt.Println(hjson.MustToString(data.Data))
			if data.Paging.Count >= int64(total) {
				break
			}
		}
		e := time.Now().UnixMilli()
		fmt.Println("elapse: ", e-s, " each ", (e-s)/int64(total))
	})
	helix.Startup()
}

func productTest() {
	_, err := product.CreateSpu(&product.SpuM{
		Collar:    collar.Build("shop", "123456"),
		Category:  119,
		Name:      "SPU-TEST-" + cast.ToString(time.Now().Unix()),
		Intro:     "SPU-INTRO",
		Brand:     "SPU-BRAND",
		Version:   0,
		MainMedia: hjson.MustToBytes(media.More{media.New(media.ImageType, "https://1.1")}),
		MoreMedia: hjson.MustToBytes(media.NewDict().Put(
			"hello",
			media.New(media.ImageType, "https://1.1").SetMeta("a", "A"),
		)),
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(spuM.ID)
	//s, err := product.CreateSpec(&product.SpuSpecSetting{
	//	Spu: spuM.ID,
	//	Specs: []*product.SpuSpec{
	//		{
	//			Spec: 999,
	//			Options: []*spec.Option{
	//				{
	//					//OptID: 0,
	//					Label: "RED SALES",
	//					Media: nil,
	//				},
	//				{
	//					//OptID: 0,
	//					Label: "BLUE SALES",
	//					Media: nil,
	//				},
	//			},
	//			//Seq: 0,
	//		},
	//		{
	//			Spec: 998,
	//			Options: []*spec.Option{
	//				{
	//					//OptID: 0,
	//					Label: "100 SALES",
	//					Media: nil,
	//				},
	//				{
	//					//OptID: 0,
	//					Label: "500 SALES",
	//					Media: nil,
	//				},
	//			},
	//			//Seq: 0,
	//		},
	//	},
	//})
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	////fmt.Println(hjson.MustToString(s))
	//var specs []spec.OptID
	//for _, spuSpecItem := range s.Specs {
	//	for _, optItem := range spuSpecItem.Options {
	//		specs = append(specs, optItem.OptID)
	//	}
	//}
	//skuSetting := &product.SkuSpecSetting{
	//	Spu:   spuM.ID,
	//	Specs: specs,
	//}
	//skuID, err := product.CreateSku(skuSetting)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	////fmt.Println(hjson.MustToString(skuSetting))
	//fmt.Println(skuID)
}
