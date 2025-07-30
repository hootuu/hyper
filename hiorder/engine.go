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
	return e.doAdvToInitial(ctx)
}
