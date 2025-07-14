package main

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/data/ctrl"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyper/hiorder"
	"github.com/hootuu/hyper/hitopup"
	_ "github.com/hootuu/hyper/hyperidx"
	"github.com/nineora/harmonic/chain"
	"github.com/nineora/harmonic/nineapi"
	"github.com/nineora/harmonic/nineora"
	"github.com/nineora/lightv/lightv"
	"github.com/spf13/cast"
	"math/rand/v2"
	"time"
)

func main() {
	helix.AfterStartup(func() {
		gjjMobi := fmt.Sprintf("13588888899")
		gjj := "gjj-" + cast.ToString(time.Now().UnixMilli())
		gjjWalletAddr, err := lightv.UserWalletHelper.Init(
			nineora.SeedOfCell(gjjMobi).ToCollar().ToSafeID(),
			gjj,
			nil,
		)
		if err != nil {
			panic(err)
		}
		fmt.Println("gjjWalletAddr:::: ", gjjWalletAddr)

		gjjAccAddr, err := lightv.UserAccHelper.GetOrInit(
			gjj,
			func() chain.Address {
				return gjjWalletAddr
			},
			func() *nineapi.Ex {
				return nil
			},
		)
		if err != nil {
			fmt.Println("UserAccHelper.GetOrInit:::: ", err)
			panic(err)
		}
		fmt.Println("gjjAccAddr:::: ", gjjAccAddr)

		topup, err := hitopup.NewTopUp(
			"XB_TUPUP",
			collar.Build(lightv.TokenLinkCode, lightv.QSTA),
			collar.Build(lightv.UserAccountLinkCode, gjj),
			func(src uint64) uint64 {
				return src
			},
		)
		if err != nil {
			panic(err)
		}

		s := time.Now()
		total := 100
		for i := 0; i < total; i++ {
			cur := i
			from := cur * 1000
			max := cur*1000 + 100
			go func() {
				for j := from; j < max; j++ {
					doHandleOrd(i, topup)
				}
			}()

		}
		elapsed := time.Since(s).Milliseconds()
		time.Sleep(100 * time.Second)
		fmt.Println()
		fmt.Println()
		fmt.Println()
		fmt.Println("=============================+>>>>>>>>>>>>>>>>>>>>>>>")
		fmt.Println(elapsed/1000, " each : ", elapsed/int64(total))
	})
	helix.Startup()
}

func doHandleOrd(seedNumb int, topup *hitopup.TopUp) {
	usrMobi := fmt.Sprintf("155%08d", seedNumb)

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
