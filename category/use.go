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

func Default() *Category {
	return gDefault
}

var gDefault *Category

func init() {
	helix.Use(helix.BuildHelix("hyper_category", func() (context.Context, error) {
		var err error
		gDefault, err = NewCategory("hyper_category_cat")
		if err != nil {
			return nil, err
		}
		return nil, nil
	}, func(ctx context.Context) {

	}))
}
