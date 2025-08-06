package hiorder

import (
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hcoin"
	"github.com/hootuu/hyle/hfsm"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyle/hypes/ex"
	"github.com/hootuu/hyper/payment"
	"github.com/hootuu/hyper/shipping"
	"gorm.io/datatypes"
	"time"
)

type OrderM struct {
	hdb.Template
	ID            ID             `gorm:"column:id;primaryKey;autoIncrement:false;"`
	Code          Code           `gorm:"column:code;index;size:32;"`
	Title         string         `gorm:"column:title;not null;size:255;"`
	Payer         collar.Link    `gorm:"column:payer;index;not null;size:128;"`
	Payee         collar.Link    `gorm:"column:payee;index;not null;size:128;"`
	Amount        hcoin.Amount   `gorm:"column:amount;autoIncrement:false;"`
	PaymentID     payment.ID     `gorm:"column:payment_id;autoIncrement:false;"`
	ShippingID    shipping.ID    `gorm:"column:shipping_id;autoIncrement:false;"`
	Status        hfsm.State     `gorm:"column:status;not null;"`
	ExStatus      hfsm.State     `gorm:"column:ex_status;not null;"`
	Matter        datatypes.JSON `gorm:"column:matter;type:jsonb;"`
	Link          collar.Link    `gorm:"column:link;index;size:128;"`
	ConsensusTime *time.Time     `gorm:"column:consensus_time;"`
	ExecutingTime *time.Time     `gorm:"column:executing_time;"`
	CanceledTime  *time.Time     `gorm:"column:canceled_time;"`
	CompletedTime *time.Time     `gorm:"column:completed_time;"`
	TimeoutTime   *time.Time     `gorm:"column:timeout_time;"`
}

func (m *OrderM) TableName() string {
	return "hyper_order"
}

func orderMto[T Matter](m *OrderM) *Order[T] {
	ord := &Order[T]{
		ID:         m.ID,
		Code:       m.Code,
		Title:      m.Title,
		Payer:      m.Payer,
		Payee:      m.Payee,
		Amount:     m.Amount,
		PaymentID:  m.PaymentID,
		ShippingID: m.ShippingID,
		Link:       m.Link,
		Status:     m.Status,
		ExStatus:   m.ExStatus,
		Ex:         ex.WithBytes(m.Ctrl, m.Tag, m.Meta),
	}
	if len(m.Matter) > 0 {
		hjson.MustOfBytes[T](&ord.Matter, m.Matter)
	}
	return ord
}
