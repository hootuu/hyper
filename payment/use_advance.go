package payment

import (
	"context"
)

func AdvPrepared(ctx context.Context, payID ID) error {
	InitIfNeeded()
	return doAdvToPrepared(ctx, payID)
}

func AdvCompleted(ctx context.Context, payID ID) error {
	InitIfNeeded()
	return doAdvToCompleted(ctx, payID)
}
