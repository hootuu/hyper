package hiorder

import "github.com/hootuu/hyper/hpay/payment"

type PaymentAltered[T Matter] struct {
	Order     *Order[T]      `json:"order"`
	PaymentID payment.ID     `json:"payment_id"`
	SrcStatus payment.Status `json:"src_status"`
	DstStatus payment.Status `json:"dst_status"`
}

func (p *PaymentAltered[T]) IsCompleted() bool {
	return p.DstStatus == payment.Completed
}
