package hiprod

import (
	"fmt"
	"github.com/hootuu/hyle/crypto/hmd5"
	"github.com/hootuu/hyle/data/hjson"
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
	vwhSku, err := vwh.DbSkuVwhGetBySku(skuM.ID)
	if err != nil {
		return nil, fmt.Errorf("vwh get error: %v", err)
	}
	return &Product{
		ID:        hmd5.MD5(fmt.Sprintf("%d:%d:%d:%d", spuM.ID, skuM.ID, 0, 0)),
		SkuID:     skuM.ID,
		SpuID:     spuM.ID,
		VwhID:     vwhSku.Vwh,
		PwhID:     vwhSku.Pwh,
		Biz:       spuM.Biz,
		Category:  spuM.Category,
		Name:      spuM.Name,
		Intro:     spuM.Intro,
		Brand:     spuM.Brand,
		Media:     *hjson.MustFromBytes[media.Dict](spuM.Media),
		Price:     vwhSku.Price,
		Inventory: vwhSku.Inventory,
	}, err
}
