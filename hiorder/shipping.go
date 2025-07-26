package hiorder

import (
	"github.com/hootuu/hyper/hshipping/shipping"
)

type ShippingAltered[T Matter] struct {
	Order      *Order[T]       `json:"order"`
	ShippingID shipping.ID     `json:"shipping_id"`
	SrcStatus  shipping.Status `json:"src_status"`
	DstStatus  shipping.Status `json:"dst_status"`
}

func (p *ShippingAltered[T]) IsCompleted() bool {
	return p.DstStatus == shipping.StatusDelivered
}
