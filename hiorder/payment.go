package hiorder

const (
	PaymentAlterTopic = "NINEPAY_PAYMENT_ALTER"
)

const (
	PaymentInitial  string = "initial"
	PaymentCanceled string = "canceled"
	PaymentTimeout  string = "timeout"
	PaymentPaid     string = "paid"
)

type PaymentPayload struct {
	OrderCollar string `json:"order_collar"`
	PaymentID   string `json:"payment_id"`
	SrcStatus   string `json:"src_status"`
	DstStatus   string `json:"dst_status"`
}

type PaymentAltered[T Matter] struct {
	Order     *Order[T] `json:"order"`
	PaymentID string    `json:"payment_id"`
	SrcStatus string    `json:"src_status"`
	DstStatus string    `json:"dst_status"`
}
