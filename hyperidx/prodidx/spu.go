package prodidx

import (
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/helix/storage/hmeili"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/hiprod/prod"
	"github.com/hootuu/hyper/hiprod/pwh"
	"github.com/hootuu/hyper/hiprod/vwh"
	"github.com/hootuu/hyper/hyperplt"
	"github.com/meilisearch/meilisearch-go"
	"go.uber.org/zap"
	"time"
)

const (
	SpuIndex      = "hyper_spu"
	spuIdxVersion = "1_0_1"
)

type TxSpuIndexer struct{}

func (idx *TxSpuIndexer) GetName() string {
	return SpuIndex
}

func (idx *TxSpuIndexer) GetVersion() string {
	return spuIdxVersion
}

func (idx *TxSpuIndexer) Setting(index meilisearch.IndexManager) error {
	filterableAttributes := []string{
		"auto_id",
		"id",
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

func (idx *TxSpuIndexer) Load(autoID int64) (hmeili.Document, error) {
	spuM, err := hdb.MustGet[prod.SpuM](hyperplt.DB(), "auto_id = ?", autoID)
	if err != nil {
		return nil, err
	}
	skuM, err := hdb.MustGet[prod.SkuM](hyperplt.DB(), "spu = ?", spuM.ID)
	if err != nil {
		return nil, err
	}

	doc := hmeili.NewMapDocument(spuM.AutoID, spuM.AutoID, spuM.UpdatedAt.UnixMilli())
	doc["sku_id"] = skuM.ID
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

	doc["media"] = spuM.Media
	doc["ctrl"] = spuM.Ctrl
	doc["tag"] = spuM.Tag
	doc["meta"] = spuM.Meta

	// fix  触发 index更新 同步spu最新数据
	go func(skuID uint64) {
		tx := hyperplt.DB()
		now := time.Now()
		_ = hdb.Update[vwh.VirtualWhSkuM](tx, map[string]interface{}{
			"updated_at": now,
		}, "sku = ?", skuID)
		_ = hdb.Update[pwh.PhysicalSkuM](tx, map[string]interface{}{
			"updated_at": now,
		}, "sku = ?", skuID)
	}(skuM.ID)

	return doc, nil
}
