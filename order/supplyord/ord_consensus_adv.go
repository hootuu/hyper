package supplyord

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hypes/ex"
	"github.com/hootuu/hyper/hiorder"
	"github.com/hootuu/hyper/hyperplt"
	"github.com/spf13/cast"
)

func doOrderConsensusAdv(ctx context.Context, orderID hiorder.ID) error {
	orderM, err := hiorder.DbMustGet(ctx, cast.ToString(orderID))
	if err != nil {
		return err
	}
	if orderM == nil {
		return nil
	}
	if orderM.Status == hiorder.Consensus {
		matter := *hjson.MustFromBytes[Matter](orderM.Matter)
		ordFactory := GetFactory()
		supOrder, err := ordFactory.Create(ctx, &CreateParas{
			Matter: matter,
			Idem:   fmt.Sprintf("ORG_PROD_ORDER_%d", orderID),
			Payer:  orderM.Payee,
			Payee:  orderM.Payer,
			Title:  orderM.Title,
			Link:   ordFactory.core.OrderCollar(orderID).Link(),
			Ex: ex.NewEx().SetMeta(map[string]interface{}{
				"orderID": orderID,
			}),
		})
		if err != nil {
			return err
		}
		err = hdb.Update[*hiorder.Order[Matter]](hyperplt.DB(), map[string]any{
			"status":         hiorder.Consensus,
			"consensus_time": orderM.ConsensusTime,
			"payment_id":     orderM.PaymentID,
			"shipping_id":    orderM.ShippingID,
		}, "id = ?", supOrder.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
