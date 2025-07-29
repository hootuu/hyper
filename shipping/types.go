package shipping

import (
	"errors"
	"github.com/hootuu/hyle/hfsm"
	"github.com/hootuu/hyle/hypes/ex"
	"github.com/hootuu/hyper/address"
	"time"
)

type ID = uint64

type Status = hfsm.State

type CourierCode = string

type Shipping struct {
	ID           ID          `json:"id"`
	Address      *Address    `json:"address"`
	CourierCode  CourierCode `json:"courier_code"`   // 快递公司编码（如：SF-顺丰，JD-京东，YTO-圆通等）
	TrackingNo   string      `json:"tracking_no"`    // 物流单号（如：SF1234567890）
	ShippedAt    *time.Time  `json:"shipped_at"`     // 发货时间（指针类型，允许为空）
	DeliveredAt  *time.Time  `json:"delivered_at"`   // 签收时间（指针类型，允许为空）
	LastSyncedAt *time.Time  `json:"last_synced_at"` // 最后一次查询物流的时间（避免频繁调用第三方）
	Status       Status      `json:"status"`         // 当前物流状态（使用预定义状态常量）
	Ex           *ex.Ex      `json:"ex"`             // 扩展信息
}

type Address struct {
	ID       address.ID        `json:"id"`
	Province string            `json:"province"`
	City     string            `json:"city"`
	District string            `json:"district"`
	Address  string            `json:"address"`
	Location *address.Location `json:"location"`
	Contact  *address.Contact  `json:"contact"`
}

func (addr *Address) Validate() error {
	if addr.Address == "" {
		return errors.New("address is empty")
	}
	if addr.Contact == nil {
		return errors.New("address is empty")
	}
	if addr.Contact.Name == "" {
		return errors.New("contact name is empty")
	}
	if addr.Contact.Mobi == "" {
		return errors.New("contact mobi is empty")
	}
	return nil
}
