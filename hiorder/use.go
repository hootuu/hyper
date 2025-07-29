package hiorder

import "github.com/hootuu/helix/helix"

var gInitialized = false

func Use() {
	helix.MustInit("hyper_order", func() error {
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
