package prodidx

import (
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/helix/storage/hmeili"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/hiprod/prod"
	"github.com/hootuu/hyper/hiprod/pwh"
	"github.com/hootuu/hyper/hyperplt"
	"github.com/meilisearch/meilisearch-go"
	"go.uber.org/zap"
)

const (
	PwhProdIndex      = "hyper_pwh_prod"
	PwhProdIdxVersion = "1_0_1"
)

type TxPwhProdIndexer struct{}

func (idx *TxPwhProdIndexer) GetName() string {
	return PwhProdIndex
}

func (idx *TxPwhProdIndexer) GetVersion() string {
	return PwhProdIdxVersion
}

func (idx *TxPwhProdIndexer) Setting(index meilisearch.IndexManager) error {
	filterableAttributes := []string{
		"auto_id",
		"id",
		"sku_id",
		"spu_id",
		"pwh_id",
		"biz",
		"name",
		"category",
		"brand",
		"spu_status",
		"created_at",
	}
	_, err := index.UpdateFilterableAttributes(&filterableAttributes)
	if err != nil {
		hlog.Err("hyper.spu.idx.Setting: Error updating filterable attributes", zap.Error(err))
		return err
	}

	sortableAttributes := []string{
		"auto_id",
		"created_at",
		"sort",
	}
	_, err = index.UpdateSortableAttributes(&sortableAttributes)
	if err != nil {
		hlog.Err("hyper.spu.Setting: Error updating sortable attributes", zap.Error(err))
		return err
	}

	searchableAttributes := []string{
		"name",
	}
	_, err = index.UpdateSettings(hmeili.BuildSearchSettings(searchableAttributes))
	if err != nil {
		hlog.Err("hyper.spu.Setting: Error updating settings", zap.Error(err))
		return err
	}
	return nil
}

func (idx *TxPwhProdIndexer) Load(autoID int64) (hmeili.Document, error) {
	pwhSkuM, err := hdb.MustGet[pwh.PhysicalSkuM](hyperplt.DB(), "auto_id = ?", autoID)
	if err != nil {
		return nil, err
	}
	skuM, err := hdb.MustGet[prod.SkuM](hyperplt.DB(), "id = ?", pwhSkuM.SKU)
	if err != nil {
		return nil, err
	}
	spuM, err := hdb.MustGet[prod.SpuM](hyperplt.DB(), "id = ?", skuM.Spu)
	if err != nil {
		return nil, err
	}

	doc := hmeili.NewMapDocument(pwhSkuM.AutoID, pwhSkuM.AutoID, pwhSkuM.UpdatedAt.UnixMilli())
	doc["sku_id"] = pwhSkuM.SKU
	doc["pwh_id"] = pwhSkuM.PWH

	doc["spu_id"] = spuM.ID
	doc["biz"] = spuM.Biz
	doc["category"] = spuM.Category
	doc["name"] = spuM.Name
	doc["intro"] = spuM.Intro
	doc["brand"] = spuM.Brand
	doc["base_price"] = spuM.Price
	doc["cost_price"] = spuM.Cost
	doc["spu_status"] = spuM.Available
	doc["created_at"] = spuM.CreatedAt.Unix()
	doc["sku_stock"] = pwhSkuM.Available
	doc["lock_stock"] = pwhSkuM.Locked

	doc["media"] = spuM.Media
	doc["ctrl"] = spuM.Ctrl
	doc["tag"] = spuM.Tag
	doc["meta"] = spuM.Meta

	return doc, nil
}
