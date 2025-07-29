package payment

import "context"

func AdvJobToPrepared(ctx context.Context, pid ID, seq int) error {
	InitIfNeeded()
	return doAdvJobToPrepared(ctx, pid, seq)
}

func AdvJobToCompleted(ctx context.Context, pid ID, seq int, payNo string) error {
	InitIfNeeded()
	return doAdvJobToCompleted(ctx, pid, seq, payNo)
}
