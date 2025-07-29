package shipping

import (
	"github.com/hootuu/helix/storage/hdb"
	"gorm.io/datatypes"
	"time"
)

type ShipM struct {
	hdb.Template
	ID               ID             `gorm:"column:id;primaryKey;autoIncrement:false;"`
	BizCode          string         `gorm:"column:biz_code;not null;size:32;"`
	BizID            string         `gorm:"column:biz_id;not null;size:128;"`
	CourierCode      CourierCode    `gorm:"column:courier_code;size:32;"`
	TrackingNo       string         `gorm:"column:tracking_no;size:64;"`
	Status           Status         `gorm:"column:status;"`
	Address          datatypes.JSON `gorm:"column:address;type:jsonb;"`
	Timeout          time.Duration  `gorm:"column:timeout;"`
	TimeoutCompleted bool           `gorm:"column:timeout_completed;"`
	SubmittedTime    *time.Time     `gorm:"column:submitted_time;"`
	FailedTime       *time.Time     `gorm:"column:failed_time;"`
	CanceledTime     *time.Time     `gorm:"column:canceled_time;"`
	CompletedTime    *time.Time     `gorm:"column:completed_time;"`
}

func (m *ShipM) TableName() string {
	return "hyper_shipping"
}
