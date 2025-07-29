package payment

import "context"

func Prepare(ctx context.Context, payID ID) error {
	InitIfNeeded()
	return doPrepare(ctx, payID)
}
