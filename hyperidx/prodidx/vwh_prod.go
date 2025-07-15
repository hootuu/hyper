package prodidx

import (
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/helix/storage/hmeili"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/hiprod/prod"
	"github.com/hootuu/hyper/hiprod/vwh"
	"github.com/hootuu/hyper/hyperplt"
	"github.com/meilisearch/meilisearch-go"
	"go.uber.org/zap"
)

const (
	VwhProdIndex      = "hyper_prod"
	vwhProdIdxVersion = "1_0_0"
)

type TxVwhProdIndexer struct{}

func (idx *TxVwhProdIndexer) GetName() string {
	return VwhProdIndex
}

func (idx *TxVwhProdIndexer) GetVersion() string {
	return vwhProdIdxVersion
}

func (idx *TxVwhProdIndexer) Setting(index meilisearch.IndexManager) error {
	filterableAttributes := []string{
		"auto_id",
		"id",
		"sku_id",
		"spu_id",
		"vwh_id",
		"pwh_id",
		"biz",
		"category",
		"brand",
	}
	_, err := index.UpdateFilterableAttributes(&filterableAttributes)
	if err != nil {
		hlog.Err("hyper.spu.idx.Setting: Error updating filterable attributes", zap.Error(err))
		return err
	}

	sortableAttributes := []string{
		"auto_id",
		"timestamp",
		"price",
	}
	_, err = index.UpdateSortableAttributes(&sortableAttributes)
	if err != nil {
		hlog.Err("hyper.spu.Setting: Error updating sortable attributes", zap.Error(err))
		return err
	}
	return nil
}

func (idx *TxVwhProdIndexer) Load(autoID int64) (hmeili.Document, error) {
	vwhSkuM, err := hdb.MustGet[vwh.VirtualWhSkuM](hyperplt.DB(), "auto_id = ?", autoID)
	if err != nil {
		return nil, err
	}
	//pwhSkuM, err := hdb.MustGet[pwh.PhysicalSkuM](hyperplt.DB(),
	//	"pwh = ? AND sku = ?", vwhSkuM)
	skuM, err := hdb.MustGet[prod.SkuM](hyperplt.DB(), "id = ?", vwhSkuM.Sku)
	if err != nil {
		return nil, err
	}
	spuM, err := hdb.MustGet[prod.SpuM](hyperplt.DB(), "id = ?", skuM.Spu)
	if err != nil {
		return nil, err
	}

	doc := hmeili.NewMapDocument(vwhSkuM.AutoID, vwhSkuM.AutoID, vwhSkuM.UpdatedAt.UnixMilli())
	doc["vwh_id"] = vwhSkuM.Vwh
	doc["sku_id"] = vwhSkuM.Sku
	doc["pwh_id"] = vwhSkuM.Pwh
	doc["price"] = vwhSkuM.Price
	doc["inventory"] = vwhSkuM.Inventory

	doc["spu_id"] = spuM.ID
	doc["biz"] = spuM.Biz
	doc["category"] = spuM.Category
	doc["name"] = spuM.Name
	doc["intro"] = spuM.Intro
	doc["brand"] = spuM.Brand

	doc["media"] = spuM.Media
	doc["ctrl"] = spuM.Ctrl
	doc["tag"] = spuM.Tag
	doc["meta"] = spuM.Meta

	return doc, nil
}
