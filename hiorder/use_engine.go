package hiorder

import "context"

func (e *Engine[T]) AdvToCompleted(ctx context.Context) error {
	return e.doAdvToCompleted(ctx)
}
