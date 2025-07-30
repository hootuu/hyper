package main

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/data/ctrl"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/data/idx"
	"github.com/hootuu/hyle/data/tag"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyper/hitopup"
	"github.com/nineora/lightv/_examples/tools"
	"github.com/nineora/lightv/qing"
	"github.com/spf13/cast"
	"time"
)

func main() {
	helix.AfterStartup(func() {
		ctx := context.WithValue(context.Background(), hlog.TraceIdKey, idx.New())

		uid, err := tools.CreateUser()
		if err != nil {
			panic(err)
		}

		xcAddr, err := qing.XC().Address(ctx)
		if err != nil {
			panic(err)
		}

		usrAccAddr, err := qing.MbrAccXC().Address(ctx, uid)
		if err != nil {
			panic(err)
		}
		fmt.Println("usrAccAddr: ", usrAccAddr)

		//创建充值订单工厂： 应该使用全局变量
		gTopup, err := hitopup.NewTopUp(
			"XB_DEPOSIT_TUPUP",
			xcAddr,
			func(src uint64) uint64 {
				return src
			},
		)
		if err != nil {
			panic(err)
		}

		ord, err := gTopup.TopUpCreate(ctx, hitopup.TopUpParas{
			Idem:           idx.New(),
			Title:          "这是一笔充值" + cast.ToString(time.Now().UnixMilli()),
			Payer:          collar.Build("user", uid).Link(),
			InAccountAddr:  usrAccAddr,
			Amount:         9860,
			PayChannelCode: "ALIPAY",
			Ctrl:           ctrl.MustNewCtrl(),
			Tag:            tag.NewTag("ALIPAY"),
			Meta:           dict.NewDict().Set("alipay_id", "16880929"),
		})
		if err != nil {
			panic(err)
		}
		fmt.Println("[ORDER_ID]: ", hjson.MustToString(ord))

		err = gTopup.TopUpPaymentPrepared(ctx, ord.ID)
		if err != nil {
			panic(err)
		}

		err = gTopup.TopUpPaymentCompleted(ctx, ord.ID, "xxx")
		if err != nil {
			panic(err)
		}

		//time.Sleep(10 * time.Second)
		//pageData, err := hyperidx.Filter(
		//	hyperidx.OrdIndex,
		//	"payer_id = '"+uid+"'",
		//	[]string{"status:desc"},
		//	pagination.PageNormal(),
		//)
		//if err != nil {
		//	panic(err)
		//}
		//fmt.Println(hjson.MustToString(pageData))
	})
	helix.Startup()
}
