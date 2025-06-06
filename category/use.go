package category

import (
	"context"
	"github.com/hootuu/helix/helix"
)

func NewCategory(code string) (*Category, error) {
	if err := helix.CheckCode(code); err != nil {
		return nil, err
	}
	cate, err := newCategory(code, 8, []uint{3, 3, 3, 3})
	if err != nil {
		return nil, err
	}
	helix.Use(cate.Helix())
	return cate, nil
}

func Frontend() *Category {
	return gFrontend
}

func Backend() *Category {
	return gBackend
}

var gFrontend *Category
var gBackend *Category

func init() {
	helix.Use(helix.BuildHelix("hyper_category", func() (context.Context, error) {
		var err error
		gFrontend, err = NewCategory("hyper_category_frontend")
		if err != nil {
			return nil, err
		}
		gBackend, err = NewCategory("hyper_category_backend")
		if err != nil {
			return nil, err
		}
		return nil, nil
	}, func(ctx context.Context) {

	}))
}
