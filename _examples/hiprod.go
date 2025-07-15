package main

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/data/ctrl"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/data/idx"
	"github.com/hootuu/hyle/data/pagination"
	"github.com/hootuu/hyle/data/tag"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/media"
	"github.com/hootuu/hyper/hiprod"
	"github.com/hootuu/hyper/hyperidx"
	"github.com/spf13/cast"
	"time"
)

func main() {
	helix.AfterStartup(func() {
		ctx := context.WithValue(context.Background(), hlog.TraceIdKey, idx.New())
		id, err := hiprod.CreateProduct(ctx, &hiprod.ProdCreateParas{
			Biz:      "QYZ",
			Category: 0,
			Brand:    0,
			Name:     "这是一个测试商品",
			Intro:    "这是一个测试商品这是一个测试商品这是一个测试商品",
			Media: media.NewDict().Put(
				"main",
				media.New(media.ImageType, "https://abc.jpg")),
			Price:     10000,
			Cost:      9000,
			Inventory: 1000,
			PwhID:     0,
			VwhID:     0,
			Ctrl:      ctrl.MustNewCtrl().MustSet(1, true),
			Tag:       tag.NewTag("HOT"),
			Meta:      dict.NewDict().Set("alias", "A|B|C|D|E"),
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(id)

		time.Sleep(3 * time.Second)
		pageData, err := hyperidx.Filter(
			hyperidx.ProductIndex,
			"sku_id = "+cast.ToString(id),
			[]string{},
			pagination.PageNormal(),
		)
		if err != nil {
			panic(err)
		}
		fmt.Println(hjson.MustToString(pageData))
	})
	helix.Startup()
}
