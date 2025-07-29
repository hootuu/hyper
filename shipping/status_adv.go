package shipping

import (
	"context"
	"errors"
	"fmt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/hfsm"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/hyperplt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
)

var gShippingFSM *hfsm.Machine
var gShippingFsmOnce sync.Once

func fsm() *hfsm.Machine {
	if gShippingFSM == nil {
		gShippingFsmOnce.Do(func() {
			gShippingFSM = hfsm.NewMachine().
				AddTransition(Initialized, SubmitEvent, advToSubmitted).
				AddTransition(Initialized, CancelEvent, advToCanceled).
				AddTransition(Submitted, CancelEvent, advToCanceled).
				AddTransition(Submitted, FailEvent, advToFailed).
				AddTransition(Submitted, CompleteEvent, advToCompleted)
		})
	}
	return gShippingFSM
}

func doAdvance(
	ctx context.Context,
	id ID,
	event hfsm.Event,
	mutSet func(shipM *ShipM, mustMut map[string]any),
) (err error) {
	if hlog.IsElapseDetail() {
		defer hlog.ElapseWithCtx(ctx, "hyper.shipping.doAdvance",
			hlog.F(zap.Uint64("id", id), zap.Int("event", int(event))),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err)}
				}
				return nil
			},
		)()
	}
	var shipM *ShipM
	shipM, err = hdb.Get[ShipM](hyperplt.Tx(ctx), "id = ?", id)
	if err != nil {
		return err
	}
	if shipM == nil {
		hlog.TraceFix("hyper.shipping.doAdvance", ctx, errors.New("no such shipping"),
			zap.Uint64("shipping_id", id))
		return nil
	}
	data := make(map[string]any)
	if mutSet != nil {
		mutSet(shipM, data)
	}
	_, err = fsm().Handle(ctx, shipM.Status, event, dict.NewDict().Set("shipping_id", id), data)
	if err != nil {
		return err
	}
	return nil
}

func doSetStatus(
	ctx context.Context,
	id ID,
	targetStatus Status,
	srcStatus Status,
	data dict.Dict,
) (err error) {
	if hlog.IsElapseFunction() {
		defer hlog.ElapseWithCtx(ctx, "hyper.shipping.doSetStatus",
			hlog.F(zap.Uint64("id", id), zap.Int("target", int(targetStatus)),
				zap.Int("src", int(srcStatus))),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err)}
				}
				return nil
			},
		)()
	}
	tx := hyperplt.Tx(ctx)
	shipM, err := hdb.Get[ShipM](hyperplt.Tx(ctx), "id = ?", id)
	if err != nil {
		return err
	}
	if shipM == nil {
		return fmt.Errorf("no such shipping: id: %d", id)
	}

	mut := data
	if mut == nil {
		mut = make(map[string]any)
	}
	mut["status"] = targetStatus
	switch targetStatus {
	case Submitted:
		mut["timeout_completed"] = gorm.Expr("CURRENT_TIMESTAMP")
	case Failed:
		mut["failed_time"] = gorm.Expr("CURRENT_TIMESTAMP")
	case Canceled:
		mut["canceled_time"] = gorm.Expr("CURRENT_TIMESTAMP")
	case Completed:
		mut["completed_time"] = gorm.Expr("CURRENT_TIMESTAMP")
	default:
	}
	rows, err := hdb.UpdateX[ShipM](tx, mut, "id = ? AND status = ?", id, srcStatus)
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("payment[id=%d, status=%d] not exist", id, srcStatus)
	}
	mqPublishShippingAlter(&AlterPayload{
		ShippingID: shipM.ID,
		BizCode:    shipM.BizCode,
		BizID:      shipM.BizID,
		Src:        srcStatus,
		Dst:        targetStatus,
	})
	return nil
}
