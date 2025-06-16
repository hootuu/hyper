package feedback

import (
	"context"
	"github.com/hootuu/helix/components/sattva"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hypes/media"
)

func Add(
	person sattva.Identification,
	title string,
	desc string,
	media media.Dict,
) (string, error) {
	fbM, err := addFeedback(&FbM{
		Person:      person,
		Title:       title,
		Description: desc,
		Media:       hjson.MustToBytes(media),
	})
	if err != nil {
		return "", err
	}
	return fbM.ID, nil
}

func init() {
	helix.Use(helix.BuildHelix("hyper_feedback", func() (context.Context, error) {
		err := zplt.HelixPgDB().PG().AutoMigrate(
			&FbM{},
		)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}, func(ctx context.Context) {

	}))
}
