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
	"github.com/spf13/cast"
)

type Order[T Matter] struct {
	ID           ID             `json:"id"`
	Code         Code           `json:"code"`
	Title        string         `json:"title"`
	Payer        collar.Collar  `json:"payer"`
	PayerAccount collar.Collar  `json:"payer_account"`
	Payee        collar.Collar  `json:"payee"`
	PayeeAccount collar.Collar  `json:"payee_account"`
	Matter       T              `json:"matter"`
	Currency     hcoin.Currency `json:"currency"`
	Amount       hcoin.Amount   `json:"amount"`
	Status       Status         `json:"status"`
	Ctrl         ctrl.Ctrl      `json:"ctrl"`
	Tag          tag.Tag        `json:"tag"`
	Meta         dict.Dict      `json:"meta"`
}

func (ord *Order[T]) toModel() *OrderM {
	m := &OrderM{
		Template: hdb.Template{
			Ctrl: ord.Ctrl,
			Tag:  hjson.MustToBytes(ord.Tag),
			Meta: hjson.MustToBytes(ord.Meta),
		},
		ID:           ord.ID,
		Code:         ord.Code,
		Title:        ord.Title,
		Payer:        ord.Payer.ToID(),
		PayerAccount: ord.PayerAccount.ToID(),
		Payee:        ord.Payee.ToID(),
		PayeeAccount: ord.PayeeAccount.ToID(),
		Currency:     ord.Currency,
		Amount:       ord.Amount,
		Status:       ord.Status,
		Matter:       hjson.MustToBytes(ord.Matter),
	}
	return m
}

func (ord *Order[T]) BuildCollar() collar.Collar {
	return collar.Build(fmt.Sprintf("HIORD_%s", ord.Code), cast.ToString(ord.ID))
}
