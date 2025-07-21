package hpay

import (
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyper/hpay/ninejob"
	"github.com/hootuu/hyper/hpay/payment"
)

func init() {
	helix.AfterStartup(func() {
		payment.RegisterJobExecutor(ninejob.NewExecutor())
	})
}
