package main

import (
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/hypes"
	"github.com/hootuu/hyper/brand"
)

func main() {
	helix.AfterStartup(func() {
		_, _ = brand.Add("BRAND-TEST",
			"BRAND-test-3",
			"BRAND-desc",
			hypes.NewMediaDict().PutAppend("main",
				hypes.NewMedia(hypes.MediaTypeVideo, "https://www.abc.com/def.jpg").SetMeta("size", "120")),
		)
	})
	helix.Startup()
}
