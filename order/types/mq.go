package ordtypes

import "github.com/hootuu/hyper/hiorder"

const (
	MqTopicProdOrderAlter = "HYPER_PROD_ORDER_ALTER"
)

type AlterPayload struct {
	OrderID hiorder.ID     `json:"order_id"`
	Code    string         `json:"code"`
	Src     hiorder.Status `json:"src"`
	Dst     hiorder.Status `json:"dst"`
}
