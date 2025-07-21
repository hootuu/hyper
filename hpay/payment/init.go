package payment

import (
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/hsys"
	"github.com/hootuu/hyper/hyperplt"
)

func init() {
	helix.AfterStartup(func() {
		err := hyperplt.DB().AutoMigrate(
			&PayM{},
			&JobM{},
		)
		if err != nil {
			hsys.Exit(err)
		}
	})
}
