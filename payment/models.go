package payment

import (
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hypes/collar"
	"gorm.io/datatypes"
	"time"
)

type PayM struct {
	hdb.Template
	ID            ID            `gorm:"column:id;primaryKey;autoIncrement:false;"`
	BizCode       string        `gorm:"column:biz_code;index;size:32;"`
	BizID         string        `gorm:"column:biz_id;index;size:128;"`
	Payer         collar.Link   `gorm:"column:payer;index;size:128;"`
	Payee         collar.Link   `gorm:"column:payee;index;size:128;"`
	Amount        uint64        `gorm:"column:amount;autoIncrement:false;"`
	Status        Status        `gorm:"column:status;"`
	Timeout       time.Duration `gorm:"column:timeout;"`
	JobCount      int           `gorm:"column:job_count;"`
	PreparedTime  *time.Time    `gorm:"column:prepared_time;"`
	ExecutingTime *time.Time    `gorm:"column:executing_time;"`
	TimeoutTime   *time.Time    `gorm:"column:timeout_time;"`
	CanceledTime  *time.Time    `gorm:"column:canceled_time;"`
	CompletedTime *time.Time    `gorm:"column:completed_time;"`
}

func (m *PayM) TableName() string {
	return "hyper_payment"
}

func (m *PayM) To() *Payment {
	return &Payment{
		ID:       m.ID,
		Payer:    m.Payer,
		Payee:    m.Payee,
		BizCode:  m.BizCode,
		BizID:    m.BizID,
		Amount:   m.Amount,
		Status:   m.Status,
		Timeout:  m.Timeout,
		JobCount: m.JobCount,
	}
}

type JobM struct {
	hdb.Basic
	ID            JobID          `gorm:"column:id;primaryKey;autoIncrement:false;"`
	Channel       Channel        `gorm:"column:channel;index;size:32;"`
	PaymentID     ID             `gorm:"column:payment_id;uniqueIndex:uk_pay_seq;autoIncrement:false;"`
	PaymentSeq    int            `gorm:"column:payment_seq;uniqueIndex:uk_pay_seq;"`
	Status        JobStatus      `gorm:"column:status;index;"`
	Timeout       time.Duration  `gorm:"column:timeout;"`
	Ctx           datatypes.JSON `gorm:"column:ctx;type:jsonb;"`
	PayNo         string         `gorm:"column:pay_no;size:100;"`
	PreparedTime  *time.Time     `gorm:"column:prepared_time;"`
	TimeoutTime   *time.Time     `gorm:"column:timeout_time;"`
	CanceledTime  *time.Time     `gorm:"column:canceled_time;"`
	CompletedTime *time.Time     `gorm:"column:completed_time;"`
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
		Timeout:    m.Timeout,
		Ctx:        *hjson.MustFromBytes[dict.Dict](m.Ctx),
	}
}
