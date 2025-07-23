package hshipping

import (
	"context"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyle/hypes/ex"
	"github.com/hootuu/hyper/hshipping/shipping"
)

type CreateParas = shipping.CreateParas

func ShippingCreate(ctx context.Context, paras CreateParas) (shipping.ID, error) {
	return shipping.Create(ctx, &paras)
}

func ShippingPrepared(
	ctx context.Context,
	uniLink collar.Link,
	courierCode shipping.CourierCode,
	trackingNo string,
	meta ex.Meta,
) (err error) {
	return shipping.Prepared(ctx, uniLink, courierCode, trackingNo, meta)
}

func ShippingCompleted(ctx context.Context, uniLink collar.Link, meta ex.Meta) (err error) {
	return shipping.Completed(ctx, uniLink, meta)
}
