package main

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/data/ctrl"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/data/pagination"
	"github.com/hootuu/hyle/data/tag"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyper/hiorder"
	"github.com/hootuu/hyper/hitopup"
	"github.com/hootuu/hyper/hyperidx"
	_ "github.com/hootuu/hyper/hyperidx"
	"github.com/nineora/lightv/qing"
	"github.com/spf13/cast"
	"time"
)

func main() {
	helix.AfterStartup(func() {
		doTopupExample()
	})
	helix.Startup()
}

func doTopupExample() {
	//保证金户头钱包
	pltXbAddr, err := qing.PltRecycleAccXB().Address(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println("pltXbAddr: ", pltXbAddr)

	//创建充值订单工厂： 应该使用全局变量
	gTopup, err := hitopup.NewTopUp(
		"XB_DEPOSIT_TUPUP",
		pltXbAddr,
		func(src uint64) uint64 {
			return src
		},
	)
	if err != nil {
		panic(err)
	}

	uid := "USER_" + cast.ToString(time.Now().UnixMilli())

	ord, err := gTopup.TopUp(context.Background(), hitopup.TopUpParas{
		Title:        "这是一笔充值" + cast.ToString(time.Now().UnixMilli()),
		Payer:        collar.Build("user", uid).Link(), // todo
		PayerAccount: collar.Build("user", uid).Link(),
		Amount:       9860,
		Ctrl:         ctrl.MustNewCtrl(),
		Tag:          tag.NewTag("ALIPAY"),
		Meta:         dict.NewDict().Set("alipay_id", "16880929"),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("<UNK>", ord)

	////模拟Ninepay发送支付成功信息
	//go func() {
	//	time.Sleep(1 * time.Second)
	//	err = zplt.HelixMqPublish(
	//		hiorder.PaymentAlterTopic,
	//		hjson.MustToBytes(hiorder.PaymentPayload{
	//			OrderCollar: ord.BuildCollar().ToID(),
	//			PaymentID:   "16880929",
	//			SrcStatus:   hiorder.PaymentInitial,
	//			DstStatus:   hiorder.PaymentPaid,
	//		}),
	//	)
	//}()

	time.Sleep(10 * time.Second)
	pageData, err := hyperidx.Filter(
		hyperidx.OrdIndex,
		"payer_id = '"+uid+"'",
		[]string{"status:desc"},
		pagination.PageNormal(),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(hjson.MustToString(pageData))
}

//
//func doHandleOrd(seedNumb int, topup *hitopup.TopUp) {
//	usrMobi := fmt.Sprintf("112%08d", seedNumb)
//
//	//qcinT, err := lightv.QstaHelper.MustGetToken()
//	//if err != nil {
//	//	panic(err)
//	//}
//	//fmt.Println(hjson.MustToString(qcinT))
//
//	uid := "user-" + cast.ToString(time.Now().UnixMilli())
//	usrWalletAddr, err := lightv.UserWalletHelper.Init(
//		nineora.SeedOfCell(usrMobi).ToCollar().ToSafeID(),
//		uid,
//		nil,
//	)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println("usrWalletAddr:::: ", usrWalletAddr)
//
//	usrAccAddr, err := lightv.UserAccHelper.GetOrInit(
//		uid,
//		func() chain.Address {
//			return usrWalletAddr
//		},
//		func() *nineapi.Ex {
//			return nil
//		},
//	)
//	if err != nil {
//		fmt.Println("UserAccHelper.GetOrInit:::: ", err)
//		panic(err)
//	}
//	fmt.Println("usrAccAddr:::: ", usrAccAddr)
//
//	s := time.Now()
//	wg := sync.WaitGroup{}
//	wg.Add(100 * 100)
//
//	for i := 0; i < 100; i++ {
//		go func() {
//			for j := 0; j < 100; j++ {
//
//				//var ord *hiorder.Order[hitopup.Matter]
//				ord, err := topup.TopUp(context.Background(), hitopup.TopUpParas{
//					Title:        "这是一笔充值" + cast.ToString(seedNumb),
//					Payer:        collar.Build("user", uid),
//					PayerAccount: collar.Build("user", uid),
//					Amount:       10001,
//					Ctrl:         ctrl.MustNewCtrl(),
//					Tag:          nil,
//					Meta:         dict.NewDict().Set("title", "<UNK>"),
//				})
//				if err != nil {
//					panic(err)
//				}
//				wg.Done()
//
//				go func() {
//					time.Sleep(time.Duration(rand.UintN(10)) * time.Second)
//					err = zplt.HelixMqPublish(
//						hiorder.PaymentAlterTopic,
//						hjson.MustToBytes(hiorder.PaymentPayload{
//							OrderCollar: string(ord.BuildCollar().ToID()),
//							PaymentID:   "111111111",
//							SrcStatus:   hiorder.PaymentInitial,
//							DstStatus:   hiorder.PaymentPaid,
//						}),
//					)
//				}()
//			}
//		}()
//	}
//	wg.Wait()
//	elapsed := time.Since(s).Milliseconds()
//
//	time.Sleep(5 * time.Minute)
//	fmt.Println()
//	fmt.Println()
//	fmt.Println()
//	fmt.Println("=============================+>>>>>>>>>>>>>>>>>>>>>>>")
//	fmt.Println(elapsed/1000, " each : ", elapsed/int64(100*100))
//
//}
