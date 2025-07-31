package payment

import (
	"github.com/hootuu/helix/helix"
)

func init() {
	RegisterJobExecutor(NewNineExecutor())
	RegisterJobExecutor(NewThirdExecutor())
	helix.Ready(func() {
		InitIfNeeded()
	})
}

var gInitialized = false

func Use() {
	helix.MustInit("hyper_payment", func() error {
		if err := doDbInit(); err != nil {
			return err
		}
		if err := doMqInit(); err != nil {
			return err
		}
		gInitialized = true
		return nil
	})
}

func InitIfNeeded() {
	if !gInitialized {
		Use()
	}
}
