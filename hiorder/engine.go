package hiorder

import (
	"context"
	"github.com/hootuu/hyle/hfsm"
)

type Engine[T Matter] struct {
	ord  *Order[T]
	deal Deal[T]
	f    *Factory[T]
	fsm  *hfsm.Machine
}

func newEngine[T Matter](deal Deal[T], ord *Order[T], f *Factory[T]) *Engine[T] {
	e := &Engine[T]{
		ord:  ord,
		deal: deal,
		f:    f,
	}
	return e
}

func (e *Engine[T]) GetOrder() *Order[T] {
	return e.ord
}

func (e *Engine[T]) Submit(ctx context.Context) error {
	return e.doAdvToInitial(ctx, e.f.nextID())
}

//
//func (e *Engine[T]) Consense(ctx context.Context) error {
//	target, err := e.fsm.Handle(ctx, e.ord.Status, ConsenseEvent)
//	if err != nil {
//		return err
//	}
//	e.ord.Status = target
//	return nil
//}
//
//func (e *Engine[T]) Execute(ctx context.Context) (err error) {
//	target, err := e.fsm.Handle(ctx, e.ord.Status, ExecuteEvent)
//	if err != nil {
//		return err
//	}
//	e.ord.Status = target
//	return nil
//}

//	func (e *Engine[T]) Complete(ctx context.Context) (err error) {
//		target, err := e.fsm.Handle(ctx, e.ord.Status, CompleteEvent)
//		if err != nil {
//			return err
//		}
//		e.ord.Status = target
//		return nil
//	}
//
//	func (e *Engine[T]) Cancel(ctx context.Context) (err error) {
//		target, err := e.fsm.Handle(ctx, e.ord.Status, CancelEvent)
//		if err != nil {
//			return err
//		}
//		e.ord.Status = target
//		return nil
//	}
//func (e *Engine[T]) onSubmit(
//	ctx context.Context,
//	_ hfsm.State,
//	_ hfsm.Event,
//	_ ...dict.Dict,
//) (target hfsm.State, err error) {
//	defer hlog.ElapseWithCtx(ctx, "hiorder.Submit", hlog.F(zap.Uint64("ord.id", e.ord.ID)),
//		func() []zap.Field {
//			if err != nil {
//				return []zap.Field{
//					zap.Any("ord", e.ord),
//					zap.Error(err),
//				}
//			}
//			return []zap.Field{
//				zap.Uint64("ord.id", e.ord.ID),
//			}
//		})()
//	//e.ord.ID = e.f.nextID() //todo recheck will delete
//	ordM := e.ord.toModel()
//
//	tx := hyperplt.Tx(ctx)
//	err = hdb.Tx(tx, func(tx *gorm.DB) error {
//		innerCtx := hdb.TxCtx(tx)
//		err = e.deal.Before(innerCtx, e.ord.Status, Initial)
//		if err != nil {
//			hlog.Err("hyper.order.Submit: dealer.Before", zap.Error(err))
//			return err
//		}
//		ordM.Status = Initial
//		err = hdb.Create[OrderM](tx, ordM)
//		if err != nil {
//			hlog.Err("hyper.order.Submit: db.Create", zap.Error(err))
//			return err
//		}
//		err = e.deal.After(innerCtx, e.ord.Status, Initial)
//		if err != nil {
//			hlog.Err("hyper.order.Submit: dealer.After", zap.Error(err))
//			return err
//		}
//		return nil
//	})
//	if err != nil {
//		hlog.Err("hyper.order.Submit: Tx", zap.Error(err))
//		return e.ord.Status, err
//	}
//
//	return Initial, nil
//}

//
//func (e *Engine[T]) onAdvanceWrap(target Status) hfsm.Transition {
//	return func(
//		ctx context.Context,
//		_ hfsm.State,
//		_ hfsm.Event,
//		_ ...dict.Dict,
//	) (hfsm.State, error) {
//		return e.onAdvance(ctx, target)
//	}
//}
//
//func (e *Engine[T]) onAdvance(
//	ctx context.Context,
//	target hfsm.State,
//) (hfsm.State, error) {
//	var err error
//	defer hlog.ElapseWithCtx(ctx, "hiorder.onAdvance",
//		hlog.F(
//			zap.Uint64("ord.id", e.ord.ID),
//			zap.Int("src", int(e.ord.Status)),
//			zap.Int("target", int(target)),
//		),
//		func() []zap.Field {
//			if err != nil {
//				return []zap.Field{
//					zap.Any("ord", e.ord),
//					zap.Error(err),
//				}
//			}
//			return []zap.Field{zap.Uint64("ord.id", e.ord.ID)}
//		})()
//	tx := hyperplt.Tx(ctx)
//	err = hdb.Tx(tx, func(tx *gorm.DB) error {
//		innerCtx := hdb.TxCtx(tx)
//		err = e.deal.Before(innerCtx, e.ord.Status, target)
//		if err != nil {
//			hlog.Err("hyper.order.onAdvance: dealer.Before", zap.Error(err))
//			return err
//		}
//		row, err := hdb.UpdateX[OrderM](tx, map[string]any{
//			"status": target,
//		}, "id = ? AND status = ?", e.ord.ID, e.ord.Status)
//		if err != nil {
//			hlog.Err("hyper.order.onAdvance: db.After", zap.Error(err))
//			return err
//		}
//		if row == 0 {
//			hlog.Err("hyper.order.onAdvance: version rebase",
//				zap.Uint64("ord.id", e.ord.ID))
//			return err
//		}
//		err = e.deal.After(innerCtx, e.ord.Status, target)
//		if err != nil {
//			hlog.Err("hyper.order.onAdvance: dealer.After", zap.Error(err))
//			return err
//		}
//		return nil
//	})
//	if err != nil {
//		hlog.Err("hyper.order.onAdvance: Tx", zap.Error(err))
//		return e.ord.Status, err
//	}
//
//	return target, nil
//}
