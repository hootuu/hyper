package main

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/data/idx"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyle/hypes/ex"
	"github.com/hootuu/hyper/hpay"
	"github.com/hootuu/hyper/hpay/ninejob"
	"github.com/hootuu/hyper/hpay/payment"
	"github.com/nineora/harmonic/harmonic"
	"github.com/nineora/lightv/_examples/tools"
	"github.com/nineora/lightv/qing"
	"time"
)

func main() {
	helix.AfterStartup(func() {
		ctx := context.WithValue(context.Background(), hlog.TraceIdKey, idx.New())
		uid := fmt.Sprintf("uid_%d", time.Now().UnixMilli())
		exM := ex.EmptyEx()
		exM.Tag.Append("xxx")
		exM.Meta.Set("xxx", "xxx")

		uid, err := tools.CreateUser()
		if err != nil {
			panic(err)
		}

		uid2, err := tools.CreateUser()
		if err != nil {
			panic(err)
		}

		xcAddr, err := qing.XC().Address(ctx)
		if err != nil {
			panic(err)
		}

		usrXcAddr, err := qing.MbrAccXC().Address(ctx, uid)
		if err != nil {
			panic(err)
		}
		fmt.Println("The user XC address is", usrXcAddr)

		usrXcAddr2, err := qing.MbrAccXC().Address(ctx, uid2)
		if err != nil {
			panic(err)
		}
		fmt.Println("The user XC address is", usrXcAddr)

		_, err = qing.XC().Mint(ctx, harmonic.TokenMintParas{
			Recipient:  usrXcAddr,
			Amount:     1000000,
			LockAmount: 0,
			Biz:        "RECHARGE",
			Ex:         nil,
			Link:       collar.Build("RE", "111111").Link(),
		})
		if err != nil {
			panic(err)
		}

		payId, err := hpay.Create(ctx, hpay.CreateParas{
			Payer:   collar.Build("USER", uid).Link(),
			Payee:   collar.Build("USER", uid).Link(),
			BizLink: collar.Build("BUY", uid).Link(),
			Amount:  900,
			Ex:      exM,
			Jobs: []payment.JobDefine{
				ninejob.Job{
					Mint:   xcAddr,
					Payer:  usrXcAddr,
					Payee:  usrXcAddr2,
					Amount: 900,
				},
			},
		})
		if err != nil {
			panic(err)
		}
		fmt.Println("PaymentID: ", payId)
		jobID := payment.BuildJobID(payId, 0)
		fmt.Println("JobID: ", jobID)
		err = hpay.Prepare(ctx, payId)
		if err != nil {
			panic(err)
		}
		err = hpay.Advance(ctx, payId)
		if err != nil {
			panic(err)
		}
	})
	helix.Startup()
}
