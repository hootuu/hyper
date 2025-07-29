package shipping

import (
	"github.com/hootuu/hyle/hsys"
	"github.com/hootuu/hyper/hyperplt"
)

func doDbInit() error {
	if hsys.RunMode().IsRd() {
		err := hyperplt.DB().AutoMigrate(
			&ShipM{},
		)
		if err != nil {
			return err
		}
	}
	return nil
}
