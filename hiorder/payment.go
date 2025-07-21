package hiorder

import "github.com/hootuu/hyper/hpay/payment"

// const (
//
//	PaymentAlterTopic = "NINEPAY_PAYMENT_ALTER"
//
// )
//
// const (
//
//	PaymentInitial  string = "initial"
//	PaymentCanceled string = "canceled"
//	PaymentTimeout  string = "timeout"
//	PaymentPaid     string = "paid"
//
// )
//
//	type PaymentPayload struct {
//		OrderCollar string `json:"order_collar"`
//		PaymentID   string `json:"payment_id"`
//		SrcStatus   string `json:"src_status"`
//		DstStatus   string `json:"dst_status"`
//	}
type PaymentAltered[T Matter] struct {
	Order     *Order[T]      `json:"order"`
	PaymentID payment.ID     `json:"payment_id"`
	SrcStatus payment.Status `json:"src_status"`
	DstStatus payment.Status `json:"dst_status"`
}

func (p *PaymentAltered[T]) IsCompleted() bool {
	return p.DstStatus == payment.Completed
}
