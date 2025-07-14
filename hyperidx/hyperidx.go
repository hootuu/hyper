package hyperidx

import (
	"github.com/hootuu/helix/storage/hmeili"
	"github.com/hootuu/hyle/data/pagination"
	"github.com/hootuu/hyper/hyperplt"
)

func Filter(idx string, filter string, sort []string, page *pagination.Page) (*pagination.Pagination[any], error) {
	return hmeili.Filter(hyperplt.Meili(), idx, filter, sort, page)
}
