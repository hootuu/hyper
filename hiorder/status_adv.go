package hiorder

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
)

func (e *Engine[T]) mustFsm() *hfsm.Machine {
	if e.fsm == nil {
		e.fsm = hfsm.NewMachine().
			AddTransition(Draft, SubmitEvent, e.advToInitial).
			AddTransition(Initial, ConsenseEvent, e.advToConsensus).
			AddTransition(Initial, CancelEvent, e.advToCanceled).
			AddTransition(Initial, TimeoutEvent, e.advToTimeout).
			AddTransition(Initial, CompleteEvent, e.advToCompleted).
			AddTransition(Consensus, ExecuteEvent, e.advToExecuting).
			AddTransition(Consensus, CompleteEvent, e.advToCompleted).
			AddTransition(Executing, CompleteEvent, e.advToCompleted)
	}
	return e.fsm
}

func (e *Engine[T]) doAdvance(
	ctx context.Context,
	event hfsm.Event,
	mutSet func(ordM *OrderM, mustMut map[string]any),
) (err error) {
	id := e.ord.ID
	if hlog.IsElapseFunction() {
		defer hlog.ElapseWithCtx(ctx, "hyper.order.doAdvance",
			hlog.F(zap.Uint64("id", id), zap.Int("event", int(event))),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err)}
				}
				return nil
			},
		)()
	}
	var ordM *OrderM
	ordM, err = hdb.Get[OrderM](hyperplt.Tx(ctx), "id = ?", id)
	if err != nil {
		return err
	}
	if ordM == nil {
		hlog.TraceFix("hyper.order.doAdvance", ctx, errors.New("no such order"),
			zap.Uint64("order_id", id))
		return nil
	}
	data := make(map[string]any)
	if mutSet != nil {
		mutSet(ordM, data)
	}
	_, err = e.mustFsm().Handle(ctx, ordM.Status, event, data)
	if err != nil {
		return err
	}
	return nil
}

func (e *Engine[T]) doSetStatus(
	ctx context.Context,
	id ID,
	targetStatus Status,
	srcStatus Status,
	data dict.Dict,
) (err error) {
	if hlog.IsElapseFunction() {
		defer hlog.ElapseWithCtx(ctx, "hyper.order.doSetStatus",
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
	mut := data
	if mut == nil {
		mut = make(map[string]any)
	}
	mut["status"] = targetStatus
	switch targetStatus {
	case Consensus:
		mut["consensus_time"] = gorm.Expr("CURRENT_TIMESTAMP")
	case Executing:
		mut["executing_time"] = gorm.Expr("CURRENT_TIMESTAMP")
	case Canceled:
		mut["canceled_time"] = gorm.Expr("CURRENT_TIMESTAMP")
	case Completed:
		mut["completed_time"] = gorm.Expr("CURRENT_TIMESTAMP")
	case Timeout:
		mut["timeout_time"] = gorm.Expr("CURRENT_TIMESTAMP")
	default:
	}
	err = hdb.Tx(tx, func(tx *gorm.DB) error {
		innerCtx := hdb.TxCtx(tx, ctx)
		err := e.deal.Before(innerCtx, srcStatus, targetStatus)
		if err != nil {
			return err
		}
		rows, err := hdb.UpdateX[OrderM](tx, mut, "id = ? AND status = ?", id, srcStatus)
		if err != nil {
			return err
		}
		if rows == 0 {
			return fmt.Errorf("payment[id=%d, status=%d] not exist", id, srcStatus)
		}
		return nil
	})

	if err != nil {
		return err
	}

	err = e.deal.After(ctx, srcStatus, targetStatus)
	if err != nil {
		return err
	}

	mqPublishOrderAlter(&AlterPayload{
		OrderID: e.ord.ID,
		Code:    e.ord.Code,
		Src:     srcStatus,
		Dst:     targetStatus,
	})
	return nil
}
