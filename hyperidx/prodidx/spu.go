package prodidx

import (
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/helix/storage/hmeili"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/hiprod/prod"
	"github.com/hootuu/hyper/hyperplt"
	"github.com/meilisearch/meilisearch-go"
	"go.uber.org/zap"
)

const (
	SpuIndex      = "hyper_spu"
	spuIdxVersion = "1_0_0"
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
	}
	_, err = index.UpdateSortableAttributes(&sortableAttributes)
	if err != nil {
		hlog.Err("hyper.spu.Setting: Error updating sortable attributes", zap.Error(err))
		return err
	}
	return nil
}

func (idx *TxSpuIndexer) Load(autoID int64) (hmeili.Document, error) {
	m, err := hdb.MustGet[prod.SpuM](hyperplt.DB(), "auto_id = ?", autoID)
	if err != nil {
		return nil, err
	}
	doc := hmeili.NewMapDocument(m.ID, m.AutoID, m.UpdatedAt.UnixMilli())
	doc["biz"] = m.Biz
	doc["category"] = m.Category
	doc["name"] = m.Name
	doc["intro"] = m.Intro
	doc["brand"] = m.Brand

	doc["media"] = m.Media
	doc["ctrl"] = m.Ctrl
	doc["tag"] = m.Tag
	doc["meta"] = m.Meta

	return doc, nil
}
