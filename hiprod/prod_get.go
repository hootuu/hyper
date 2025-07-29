package hiprod

import (
	"fmt"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/data/idx"
	"github.com/hootuu/hyle/hypes/media"
	"github.com/hootuu/hyper/hiprod/prod"
	"github.com/hootuu/hyper/hiprod/vwh"
)

type ProductGetArgs struct {
	SkuID prod.SkuID
	VwhID vwh.ID
}

func ProductMustGet(args ProductGetArgs) (*Product, error) {
	skuM, err := prod.DbSkuGet(args.SkuID)
	if err != nil {
		return nil, err
	}
	if skuM == nil {
		return nil, fmt.Errorf("sku not found: %d", args.SkuID)
	}
	spuM, err := prod.DbSpuGet(skuM.Spu)
	if err != nil {
		return nil, err
	}
	if spuM == nil {
		return nil, fmt.Errorf("spu not found: %d", args.SkuID)
	}
	return &Product{
		ID:        idx.New(),
		SkuID:     skuM.ID,
		SpuID:     spuM.ID,
		VwhID:     0,
		PwhID:     0,
		Biz:       spuM.Biz,
		Category:  spuM.Category,
		Name:      spuM.Name,
		Intro:     spuM.Intro,
		Brand:     spuM.Brand,
		Media:     *hjson.MustFromBytes[media.Dict](spuM.Media),
		Price:     0, //todo
		Inventory: 0, //todo
	}, err
}
