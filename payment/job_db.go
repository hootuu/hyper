package payment

import (
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyper/hyperplt"
	"gorm.io/gorm"
)

func DbFindJobByPayID(payID ID) ([]*JobM, error) {
	return hdb.Find[JobM](func() *gorm.DB {
		return hyperplt.DB().Where("payment_id = ?", payID)
	})
}
