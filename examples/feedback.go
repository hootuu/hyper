package main

import (
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/hypes"
	"github.com/hootuu/hyper/feedback"
)

func main() {
	helix.AfterStartup(func() {
		_, _ = feedback.Add("abcd", "feedback-test", "", nil)
		_, _ = feedback.Add("abcd",
			"feedback-test-2",
			"feedback-desc",
			hypes.NewMediaDict().PutAppend("main", &hypes.Media{
				Type: hypes.MediaTypeVideo,
				Link: "https://www.abc.com/def.jpg",
			}),
		)
		_, _ = feedback.Add("abcd",
			"feedback-test-3",
			"feedback-desc",
			hypes.NewMediaDict().PutAppend("main",
				hypes.NewMedia(hypes.MediaTypeVideo, "https://www.abc.com/def.jpg").SetMeta("size", "120")),
		)
	})
	helix.Startup()
}
