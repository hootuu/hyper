package payment

import (
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hypes/collar"
	"gorm.io/datatypes"
)

type PayM struct {
	hdb.Template
	ID                ID          `gorm:"column:id;primaryKey;autoIncrement:false;"`
	BizLink           collar.Link `gorm:"column:biz_link;uniqueIndex;size:128;"`
	Payer             collar.Link `gorm:"column:payer;index;size:128;"`
	Payee             collar.Link `gorm:"column:payee;index;size:128;"`
	Amount            uint64      `gorm:"column:amount;autoIncrement:false;"`
	Status            Status      `gorm:"column:status;"`
	JobCount          int         `gorm:"column:job_count;"`
	PreparedJobCount  int         `gorm:"column:prepared_job_count;"`
	CompletedJobCount int         `gorm:"column:completed_job_count;"`
}

func (m *PayM) TableName() string {
	return "hyper_payment"
}

func (m *PayM) To() *Payment {
	return &Payment{
		ID:       m.ID,
		Payer:    m.Payer,
		Payee:    m.Payee,
		Biz:      m.BizLink,
		Amount:   m.Amount,
		Status:   m.Status,
		JobCount: m.JobCount,
	}
}

type JobM struct {
	hdb.Basic
	ID         JobID          `gorm:"column:id;primaryKey;autoIncrement:false;"`
	Channel    Channel        `gorm:"column:channel;index;size:32;"`
	PaymentID  ID             `gorm:"column:payment_id;uniqueIndex:uk_pay_seq;autoIncrement:false;"`
	PaymentSeq int            `gorm:"column:payment_seq;uniqueIndex:uk_pay_seq;"`
	Status     JobStatus      `gorm:"column:status;index;"`
	Ctx        datatypes.JSON `gorm:"column:ctx;type:jsonb;"`
	CheckCode  string         `gorm:"column:check_code;size:32;"`
}

func (m *JobM) TableName() string {
	return "hyper_payment_job"
}

func (m *JobM) To() *Job {
	return &Job{
		JobID:      m.ID,
		Channel:    m.Channel,
		PaymentID:  m.PaymentID,
		PaymentSeq: m.PaymentSeq,
		Status:     m.Status,
		Ctx:        *hjson.MustFromBytes[dict.Dict](m.Ctx),
	}
}
