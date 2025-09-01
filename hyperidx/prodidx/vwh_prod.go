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
	vwhProdIdxVersion = "1_0_1"
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
		"name",
		"category",
		"brand",
		"spu_status",
		"created_at",
		"channel",
		"available",
	}
	_, err := index.UpdateFilterableAttributes(&filterableAttributes)
	if err != nil {
		hlog.Err("hyper.spu.idx.Setting: Error updating filterable attributes", zap.Error(err))
		return err
	}

	sortableAttributes := []string{
		"auto_id",
		"created_at",
		"price",
		"available",
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

func (idx *TxVwhProdIndexer) Load(autoID int64) (hmeili.Document, error) {
	vwhSkuM, err := hdb.MustGet[vwh.VirtualWhSkuM](hyperplt.DB(), "auto_id = ?", autoID)
	if err != nil {
		return nil, err
	}
	vwhSkuExtM, err := hdb.Get[vwh.VirtualWhSkuExtM](hyperplt.DB(), "vwh = ? AND sku = ? AND pwh = ?", vwhSkuM.Vwh, vwhSkuM.Sku, vwhSkuM.Pwh)
	if err != nil {
		return nil, err
	}
	//pwhSkuM, err := hdb.MustGet[pwh.PhysicalSkuM](hyperplt.DB(),
	//	"pwh = ? AND sku = ?", vwhSkuM)
	//if err != nil {
	//	return nil, err
	//}
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
	doc["cur_stock"] = vwhSkuM.Inventory

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
	//doc["sku_stock"] = pwhSkuM.Available

	doc["media"] = spuM.Media
	doc["ctrl"] = spuM.Ctrl
	doc["tag"] = spuM.Tag
	doc["meta"] = spuM.Meta

	if vwhSkuExtM != nil {
		doc["channel"] = vwhSkuExtM.Channel
		doc["available"] = vwhSkuExtM.Available
		doc["sort"] = vwhSkuExtM.Sort
		doc["sku_ext"] = vwhSkuExtM.Meta
	} else {
		doc["channel"] = 0
		doc["available"] = false
		doc["sort"] = 0
	}

	return doc, nil
}
