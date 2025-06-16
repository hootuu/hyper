package main

import (
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/hypes/media"
	"github.com/hootuu/hyper/brand"
)

func main() {
	helix.AfterStartup(func() {
		_, _ = brand.Add("BRAND-TEST",
			"BRAND-test-3",
			"BRAND-desc",
			media.NewDict().Put("main",
				media.New(media.ImageType, "https://www.abc.com/def.jpg").SetMeta("size", "120")),
		)
	})
	helix.Startup()
}
