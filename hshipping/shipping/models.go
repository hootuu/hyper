package shipping

import (
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hypes/collar"
	"gorm.io/datatypes"
	"time"
)

type ShipM struct {
	hdb.Template
	ID           ID             `gorm:"column:id;primaryKey;autoIncrement:false;"`
	UniLink      collar.Link    `gorm:"column:uni_link;uniqueIndex;size:128;"`
	CourierCode  CourierCode    `gorm:"column:courier_code;size:32;"`
	TrackingNo   string         `gorm:"column:tracking_no;size:64;"`
	ShippedAt    *time.Time     `gorm:"column:shipped_at;"`
	LastSyncedAt *time.Time     `gorm:"column:last_synced_at;"`
	Status       Status         `gorm:"column:status;"`
	StatusText   string         `gorm:"column:status_text;size:300;"`
	Address      datatypes.JSON `gorm:"column:address;type:jsonb;"`
}

func (m *ShipM) TableName() string {
	return "hyper_shipping"
}
