package hiorder

import "context"

func (e *Engine[T]) AdvToCompleted(ctx context.Context) error {
	return e.doAdvToCompleted(ctx)
}

func (e *Engine[T]) AdvToCanceled(ctx context.Context, mutSet func(ordM *OrderM, mustMut map[string]any)) error {
	return e.doAdvToCanceled(ctx, mutSet)
}
