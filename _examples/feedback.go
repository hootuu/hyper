package main

import (
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/hypes/media"
	"github.com/hootuu/hyper/feedback"
)

func main() {
	helix.AfterStartup(func() {
		_, _ = feedback.Add("abcd", "feedback-test", "", nil)
		_, _ = feedback.Add("abcd",
			"feedback-test-2",
			"feedback-desc",
			media.NewDict().Put("main", &media.Media{
				Type: media.ImageType,
				Link: "https://www.abc.com/def.jpg",
			}),
		)
		_, _ = feedback.Add("abcd",
			"feedback-test-3",
			"feedback-desc",
			media.NewDict().Put("main", &media.Media{
				Type: media.ImageType,
				Link: "https://www.abc.com/def.jpg",
			}),
		)
	})
	helix.Startup()
}
