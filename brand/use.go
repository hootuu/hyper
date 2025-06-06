package brand

import (
	"context"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hypes/media"
)

func Add(
	name string,
	intro string,
	desc string,
	media media.Dict,
) (string, error) {
	brM, err := addBrand(&BrM{
		Name:        name,
		Intro:       intro,
		Description: desc,
		Media:       hjson.MustToBytes(media),
	})
	if err != nil {
		return "", err
	}
	return brM.ID, nil
}

func init() {
	helix.Use(helix.BuildHelix("hyper_brand", func() (context.Context, error) {
		err := zplt.HelixPgDB().PG().AutoMigrate(
			&BrM{},
		)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}, func(ctx context.Context) {

	}))
}
