package main

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/data/idx"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyle/hypes/ex"
	"github.com/hootuu/hyper/payment"
	"github.com/nineora/lightv/_examples/tools"
	"time"
)

func main() {
	helix.AfterStartup(func() {
		payment.ListeningAlter("PAY_BIZ", func(ctx context.Context, payload *payment.AlterPayload) error {
			fmt.Println("ON PAYMENT ALTER", hjson.MustToString(payload))
			return nil
		})

		ctx := context.WithValue(context.Background(), hlog.TraceIdKey, idx.New())
		uid := fmt.Sprintf("uid_%d", time.Now().UnixMilli())
		exM := ex.EmptyEx()
		exM.Tag.Append("xxx")
		exM.Meta.Set("xxx", "xxx")

		uid, err := tools.CreateUser()
		if err != nil {
			panic(err)
		}
		//
		//uid2, err := tools.CreateUser()
		//if err != nil {
		//	panic(err)
		//}

		//xcAddr, err := qing.XC().Address(ctx)
		//if err != nil {
		//	panic(err)
		//}
		//
		//usrXcAddr, err := qing.MbrAccXC().Address(ctx, uid)
		//if err != nil {
		//	panic(err)
		//}
		//fmt.Println("The user XC address is", usrXcAddr)
		//
		//usrXcAddr2, err := qing.MbrAccXC().Address(ctx, uid2)
		//if err != nil {
		//	panic(err)
		//}
		//fmt.Println("The user XC address is", usrXcAddr)
		//
		//_, err = qing.XC().Mint(ctx, harmonic.TokenMintParas{
		//	Recipient:  usrXcAddr,
		//	Amount:     1000000,
		//	LockAmount: 0,
		//	Biz:        "RECHARGE",
		//	Ex:         nil,
		//	Link:       collar.Build("RE", "111111").Link(),
		//})
		//if err != nil {
		//	panic(err)
		//}

		payId, err := payment.Create(ctx, &payment.CreateParas{
			Idem:    idx.New(),
			Payer:   collar.Build("USER", uid).Link(),
			Payee:   collar.Build("USER", uid).Link(),
			BizCode: "PAY_BIZ",
			BizID:   idx.New(),
			Amount:  900,
			Timeout: 10 * time.Second,
			Ex:      exM,
			Jobs: []payment.JobDefine{
				//&ninejob.Job{
				//	Mint:   xcAddr,
				//	Payer:  usrXcAddr,
				//	Payee:  usrXcAddr2,
				//	Amount: 800,
				//},
				&payment.ThirdJob{
					ThirdCode: "WECHAT",
					Amount:    900,
					Ex:        ex.EmptyEx(),
				},
			},
		})
		if err != nil {
			panic(err)
		}
		fmt.Println("PaymentID: ", payId)

		err = payment.Prepare(ctx, payId)
		if err != nil {
			panic(err)
		}

		err = payment.AdvJobToPrepared(ctx, payId, 1)
		if err != nil {
			panic(err)
		}

		//err = hpay.Advance(ctx, payId)
		//if err != nil {
		//	panic(err)
		//}

		//err = hpay.JobCompleted(ctx, payId, 2, "xxxx")
		//if err != nil {
		//	panic(err)
		//}
	})
	helix.Startup()
}
