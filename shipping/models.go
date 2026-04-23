package shipping

import (
	"time"

	"github.com/hootuu/helix/storage/hdb"
	"gorm.io/datatypes"
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

type ShipPkgM struct {
	hdb.Template
	ID          ID          `gorm:"column:id;primaryKey;autoIncrement:false;"`
	ShippingID  ID          `gorm:"column:shipping_id;index:uk_shipping_pkg_seq;autoIncrement:false;"`
	BizCode     string      `gorm:"column:biz_code;index:idx_shipping_pkg_biz;size:32;"`
	BizID       string      `gorm:"column:biz_id;index:idx_shipping_pkg_biz;index:idx_shipping_pkg_biz_id;size:128;"`
	PackageSeq  int         `gorm:"column:package_seq;uniqueIndex:uk_shipping_pkg_seq;"`
	CourierCode CourierCode `gorm:"column:courier_code;size:32;"`
	TrackingNo  string      `gorm:"column:tracking_no;size:64;"`
	IsPrimary   bool        `gorm:"column:is_primary;"`
}

func (m *ShipPkgM) TableName() string {
	return "hyper_shipping_package"
}
