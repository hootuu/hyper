package hiorder

import (
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/data/tag"
	"github.com/hootuu/hyle/hcoin"
	"github.com/hootuu/hyle/hfsm"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyper/hpay/payment"
	"gorm.io/datatypes"
)

type OrderM struct {
	hdb.Template
	ID        ID             `gorm:"column:id;primaryKey;autoIncrement:false;"`
	Code      Code           `gorm:"column:code;index;size:32;"`
	Title     string         `gorm:"column:title;not null;size:255;"`
	Payer     collar.Link    `gorm:"column:payer;index;not null;size:128;"`
	Payee     collar.Link    `gorm:"column:payee;index;not null;size:128;"`
	Amount    hcoin.Amount   `gorm:"column:amount;autoIncrement:false;"`
	PaymentID payment.ID     `gorm:"column:payment_id;uniqueIndex;autoIncrement:false;"`
	Status    hfsm.State     `gorm:"column:status;not null;"`
	Matter    datatypes.JSON `gorm:"column:matter;type:jsonb;"`
	UniLink   collar.Link    `gorm:"column:uni_link;index;not null;size:128;"`
}

func (m *OrderM) TableName() string {
	return "hyper_order"
}

func orderMto[T Matter](m *OrderM) *Order[T] {
	ord := &Order[T]{
		ID:        m.ID,
		Code:      m.Code,
		Title:     m.Title,
		Payer:     m.Payer,
		Payee:     m.Payee,
		Amount:    m.Amount,
		PaymentID: m.PaymentID,
		Status:    m.Status,
		Ctrl:      m.Ctrl,
		Tag:       *hjson.MustFromBytes[tag.Tag](m.Tag),
		Meta:      *hjson.MustFromBytes[dict.Dict](m.Meta),
		Matter:    *hjson.MustFromBytes[T](m.Matter),
		UniLink:   m.UniLink,
	}
	return ord
}
