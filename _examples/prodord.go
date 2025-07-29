package main

import (
	"github.com/hootuu/helix/helix"
)

func main() {
	helix.Ready(func() {
		//ord := GetID()
		//f := prodord.NewFactory(rod.Co)
		//f.Core().RegisterAlterHandle(func(ctx context.Context, payload *hiorder.AlterPayload) error {
		//	switch payload.Dst {
		//	case hiorder.Timeout:
		//		fmt.Println("do timeout this order")
		//	case hiorder.Completed:
		//		fmt.Println("do completed this order")
		//	default:
		//	}
		//	return nil
		//})
		//ord, err := f.Create(context.Background(), &prodord.CreateParas{
		//	Idem:     idx.New(),
		//	Title:    "ABC",
		//	ProdID:   "123",
		//	Payer:    "123",
		//	Payee:    "123",
		//	Quantity: 1,
		//	Amount:   120,
		//	Ex:       nil,
		//})
		//if err != nil {
		//	panic(err)
		//}
		//fmt.Println(ord)
		//eng, err := f.Core().Load(context.Background(), ord.ID)
		//if err != nil {
		//	panic(err)
		//}
		//paymentId, err := eng.SetPayment(context.Background(), &hiorder.SetPaymentParas{
		//	Idem:    idx.New(),
		//	OrderID: ord.ID,
		//	Payments: []payment.JobDefine{
		//		&payment.ThirdJob{
		//			ThirdCode: "WECHAT",
		//			Amount:    120,
		//			Ex:        ex.EmptyEx(),
		//		},
		//	},
		//	Timeout: 10 * time.Second,
		//	Ex:      nil,
		//})
		//if err != nil {
		//	panic(err)
		//}
		//err = payment.Prepare(context.Background(), paymentId)
		//if err != nil {
		//	panic(err)
		//}
		//err = payment.AdvJobToPrepared(context.Background(), paymentId, 1)
		//if err != nil {
		//	panic(err)
		//}
		////err = payment.AdvJobToCompleted(context.Background(), paymentId, 1, "ali001")
		////if err != nil {
		////	panic(err)
		////}
		////err = payment.AdvJobToCompleted(context.Background(), paymentId, 1, "ali001")
		////if err != nil {
		////	panic(err)
		////}
		//ship, err := shipping.Create(ctx, &shipping.CreateParas{
		//	Idem:    "",
		//	BizCode: "",
		//	BizID:   "",
		//	Address: nil,
		//	Ex:      nil,
		//	Time
		//})
		////shippingID, err := eng.SetShipping(dc)
		//
		//ord.Shipping
		//
		//shipping.AdvCompleted()
	})
	helix.Startup()
}
