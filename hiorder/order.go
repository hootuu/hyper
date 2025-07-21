package hiorder

import (
	"fmt"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/ctrl"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/data/tag"
	"github.com/hootuu/hyle/hcoin"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyper/hpay/payment"
	"github.com/spf13/cast"
)

type Order[T Matter] struct {
	ID        ID           `json:"id"`
	Code      Code         `json:"code"`
	Title     string       `json:"title"`
	Payer     collar.Link  `json:"payer"`
	Payee     collar.Link  `json:"payee"`
	Matter    T            `json:"matter"`
	Amount    hcoin.Amount `json:"amount"`
	PaymentID payment.ID   `json:"payment_id"`
	Status    Status       `json:"status"`
	Ctrl      ctrl.Ctrl    `json:"ctrl"`
	Tag       tag.Tag      `json:"tag"`
	Meta      dict.Dict    `json:"meta"`
}

func (ord *Order[T]) toModel() *OrderM {
	m := &OrderM{
		Template: hdb.Template{
			Ctrl: ord.Ctrl,
			Tag:  hjson.MustToBytes(ord.Tag),
			Meta: hjson.MustToBytes(ord.Meta),
		},
		ID:        ord.ID,
		Code:      ord.Code,
		Title:     ord.Title,
		Payer:     ord.Payer,
		Payee:     ord.Payee,
		Amount:    ord.Amount,
		PaymentID: ord.PaymentID,
		Status:    ord.Status,
		Matter:    hjson.MustToBytes(ord.Matter),
	}
	return m
}

func (ord *Order[T]) BuildCollar() collar.Collar {
	return collar.Build(fmt.Sprintf("HIORD_%s", ord.Code), cast.ToString(ord.ID))
}
