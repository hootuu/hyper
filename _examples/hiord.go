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
	"github.com/nineora/harmonic/chain"
	"github.com/nineora/harmonic/nineapi"
	"github.com/nineora/harmonic/nineora"
	"github.com/nineora/lightv/lightv"
	"github.com/spf13/cast"
	"math/rand/v2"
	"sync"
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
	xbDepositEmail := fmt.Sprintf("xb.deposit@lightv.com")
	xbDepositID := "xb.deposit@lightv.com"
	xdDepositWalletAddr, err := lightv.UserWalletHelper.Init(
		nineora.SeedOfEmail(xbDepositEmail).ToCollar().ToSafeID(),
		xbDepositID,
		nil,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("星币保证金钱包地址:::: ", xdDepositWalletAddr)

	xbDepositAccAddr, err := lightv.UserAccHelper.GetOrInit(
		xbDepositID,
		func() chain.Address {
			return xdDepositWalletAddr
		},
		func() *nineapi.Ex {
			return nil
		},
	)
	if err != nil {
		fmt.Println("UserAccHelper.GetOrInit:::: ", err)
		panic(err)
	}
	fmt.Println("星币保证金账户地址:::: ", xbDepositAccAddr)

	//创建充值订单工厂： 应该使用全局变量
	gTopup, err := hitopup.NewTopUp(
		"XB_DEPOSIT_TUPUP",
		collar.Build(lightv.TokenLinkCode, lightv.QSTA),
		collar.Build(lightv.UserAccountLinkCode, xbDepositID),
		func(src uint64) uint64 {
			return src
		},
	)
	if err != nil {
		panic(err)
	}

	uid := "USER_" + cast.ToString(time.Now().UnixMilli())
	//以下在CASH充值可以不用
	//usrMobi := "8618088889999"
	//usrWalletAddr, err := lightv.UserWalletHelper.Init(
	//	nineora.SeedOfCell(usrMobi).ToCollar().ToSafeID(),
	//	uid,
	//	nil,
	//)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("测试用户钱包地址:::: ", usrWalletAddr)
	//
	//usrAccAddr, err := lightv.UserAccHelper.GetOrInit(
	//	uid,
	//	func() chain.Address {
	//		return usrWalletAddr
	//	},
	//	func() *nineapi.Ex {
	//		return nil
	//	},
	//)
	//if err != nil {
	//	fmt.Println("UserAccHelper.GetOrInit:::: ", err)
	//	panic(err)
	//}
	//fmt.Println("测试用户账户地址:::: ", usrAccAddr)

	ord, err := gTopup.TopUp(context.Background(), hitopup.TopUpParas{
		Title:        "这是一笔充值" + cast.ToString(time.Now().UnixMilli()),
		Payer:        collar.Build("user", uid),
		PayerAccount: collar.Build("user", uid),
		Amount:       9860,
		Ctrl:         ctrl.MustNewCtrl(),
		Tag:          tag.NewTag("ALIPAY"),
		Meta:         dict.NewDict().Set("alipay_id", "16880929"),
	})
	if err != nil {
		panic(err)
	}

	//模拟Ninepay发送支付成功信息
	go func() {
		time.Sleep(1 * time.Second)
		err = zplt.HelixMqPublish(
			hiorder.PaymentAlterTopic,
			hjson.MustToBytes(hiorder.PaymentPayload{
				OrderCollar: ord.BuildCollar().ToID(),
				PaymentID:   "16880929",
				SrcStatus:   hiorder.PaymentInitial,
				DstStatus:   hiorder.PaymentPaid,
			}),
		)
	}()

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

func doHandleOrd(seedNumb int, topup *hitopup.TopUp) {
	usrMobi := fmt.Sprintf("112%08d", seedNumb)

	//qcinT, err := lightv.QstaHelper.MustGetToken()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(hjson.MustToString(qcinT))

	uid := "user-" + cast.ToString(time.Now().UnixMilli())
	usrWalletAddr, err := lightv.UserWalletHelper.Init(
		nineora.SeedOfCell(usrMobi).ToCollar().ToSafeID(),
		uid,
		nil,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("usrWalletAddr:::: ", usrWalletAddr)

	usrAccAddr, err := lightv.UserAccHelper.GetOrInit(
		uid,
		func() chain.Address {
			return usrWalletAddr
		},
		func() *nineapi.Ex {
			return nil
		},
	)
	if err != nil {
		fmt.Println("UserAccHelper.GetOrInit:::: ", err)
		panic(err)
	}
	fmt.Println("usrAccAddr:::: ", usrAccAddr)

	s := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(100 * 100)

	for i := 0; i < 100; i++ {
		go func() {
			for j := 0; j < 100; j++ {

				//var ord *hiorder.Order[hitopup.Matter]
				ord, err := topup.TopUp(context.Background(), hitopup.TopUpParas{
					Title:        "这是一笔充值" + cast.ToString(seedNumb),
					Payer:        collar.Build("user", uid),
					PayerAccount: collar.Build("user", uid),
					Amount:       10001,
					Ctrl:         ctrl.MustNewCtrl(),
					Tag:          nil,
					Meta:         dict.NewDict().Set("title", "<UNK>"),
				})
				if err != nil {
					panic(err)
				}
				wg.Done()

				go func() {
					time.Sleep(time.Duration(rand.UintN(10)) * time.Second)
					err = zplt.HelixMqPublish(
						hiorder.PaymentAlterTopic,
						hjson.MustToBytes(hiorder.PaymentPayload{
							OrderCollar: string(ord.BuildCollar().ToID()),
							PaymentID:   "111111111",
							SrcStatus:   hiorder.PaymentInitial,
							DstStatus:   hiorder.PaymentPaid,
						}),
					)
				}()
			}
		}()
	}
	wg.Wait()
	elapsed := time.Since(s).Milliseconds()

	time.Sleep(5 * time.Minute)
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println("=============================+>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Println(elapsed/1000, " each : ", elapsed/int64(100*100))

}
