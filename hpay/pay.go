package hpay

import (
	"context"
	"github.com/hootuu/hyper/hpay/payment"
)

type CreateParas = payment.CreateParas

func Create(ctx context.Context, paras CreateParas) (payment.ID, error) {
	return payment.Create(ctx, &paras)
}

func Prepare(ctx context.Context, id payment.ID) error {
	return payment.Prepare(ctx, id)
}

func Advance(ctx context.Context, id payment.ID) error {
	return payment.Advance(ctx, id)
}

func JobPrepared(ctx context.Context, id payment.ID, seq int, checkCode string) error {
	return payment.DoJobPrepared(ctx, id, seq, checkCode)
}

func JobCompleted(
	ctx context.Context,
	id payment.ID,
	seq int,
	checkCode string,
	payNumber string,
) error {
	return payment.DoJobCompleted(ctx, id, seq, checkCode, payNumber)
}
