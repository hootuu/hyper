package hiorder

import (
	"fmt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hcoin"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyle/hypes/ex"
	"github.com/hootuu/hyper/payment"
	"github.com/hootuu/hyper/shipping"
	"github.com/spf13/cast"
	"strings"
)

type Order[T Matter] struct {
	ID         ID           `json:"id"`
	Code       Code         `json:"code"`
	Title      string       `json:"title"`
	Payer      collar.Link  `json:"payer"`
	Payee      collar.Link  `json:"payee"`
	Matter     T            `json:"matter"`
	Amount     hcoin.Amount `json:"amount"`
	PaymentID  payment.ID   `json:"payment_id"`
	ShippingID shipping.ID  `json:"shipping_id"`
	Status     Status       `json:"status"`
	ExStatus   ExStatus     `json:"ex_status"`
	Ex         *ex.Ex       `json:"ex"`
}

func (ord *Order[T]) toModel() *OrderM {
	m := &OrderM{
		Template:   hdb.TemplateFromEx(ord.Ex),
		ID:         ord.ID,
		Code:       ord.Code,
		Title:      ord.Title,
		Payer:      ord.Payer,
		Payee:      ord.Payee,
		Amount:     ord.Amount,
		PaymentID:  ord.PaymentID,
		ShippingID: ord.ShippingID,
		Status:     ord.Status,
		ExStatus:   ord.ExStatus,
		Matter:     hjson.MustToBytes(ord.Matter),
	}
	return m
}

func (ord *Order[T]) BuildCollar() collar.Collar {
	return collar.Build(fmt.Sprintf("HIORD_%s", strings.ToUpper(ord.Code)), cast.ToString(ord.ID))
}

func (ord *Order[T]) GetDigest() ex.Meta {
	return ex.Meta{
		"code":        ord.Code,
		"id":          ord.ID,
		"payer":       ord.Payer,
		"payee":       ord.Payee,
		"amount":      ord.Amount,
		"payment_id":  ord.PaymentID,
		"shipping_id": ord.ShippingID,
		"matter":      ord.Matter.GetDigest(),
	}
}
