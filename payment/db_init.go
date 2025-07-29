package payment

import (
	"github.com/hootuu/hyle/hsys"
	"github.com/hootuu/hyper/hyperplt"
)

func doDbInit() error {
	if hsys.RunMode().IsRd() {
		err := hyperplt.DB().AutoMigrate(
			&PayM{},
			&JobM{},
		)
		if err != nil {
			return err
		}
	}
	return nil
}
