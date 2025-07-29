package payment

import (
	"github.com/hootuu/helix/helix"
)

func init() {
	RegisterJobExecutor(NewNineExecutor())
	RegisterJobExecutor(NewThirdExecutor())
}

var gInitialized = false

func Use() {
	helix.MustInit("hyper_shipping", func() error {
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
