package hyperidx

import (
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/helix/storage/hmeili"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/hyperplt"
	"github.com/hootuu/hyper/product"
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
		"category",
		"brand",
	}
	_, err := index.UpdateFilterableAttributes(&filterableAttributes)
	if err != nil {
		hlog.Err("nineidx.Setting: Error updating filterable attributes", zap.Error(err))
		return err
	}

	sortableAttributes := []string{
		"auto_id",
		"timestamp",
	}
	_, err = index.UpdateSortableAttributes(&sortableAttributes)
	if err != nil {
		hlog.Err("nineidx.Setting: Error updating sortable attributes", zap.Error(err))
		return err
	}
	return nil
}

/** TODO
Collar    collar.Collar  `gorm:"column:collar;index;size:64;"`
ID        SpuID          `gorm:"column:id;primaryKey;size:32;"`
Category  category.ID    `gorm:"column:category;index;"`
Name      string         `gorm:"column:name;size:100;"`
Intro     string         `gorm:"column:intro;size:1000;"`
Brand     brand.ID       `gorm:"column:brand;size:32;"`
Version   hdb.Version    `gorm:"column:version;"`
MainMedia datatypes.JSON `gorm:"column:main_media;type:jsonb;"` //media.More
MoreMedia datatypes.JSON `gorm:"column:more_media;type:jsonb;"` //media.Dict
*/

func (idx *TxSpuIndexer) Load(autoID int64) (hmeili.Document, error) {
	m, err := hdb.MustGet[product.SpuM](hyperplt.DB(), "auto_id = ?", autoID)
	if err != nil {
		return nil, err
	}
	doc := hmeili.NewMapDocument(m.ID, m.AutoID, m.UpdatedAt.UnixMilli())
	doc["category"] = m.Category
	doc["name"] = m.Name
	doc["intro"] = m.Intro
	doc["brand"] = m.Brand

	doc["main_media"] = m.MainMedia
	doc["more_media"] = m.MoreMedia

	//doc["more_media"] = *hjson.MustFromBytes[media.Dict](m.MoreMedia)
	//
	//if err := mixNetwork(txM.Network, func(data dict.Dict) {
	//	doc.Mix("network", data)
	//}); err != nil {
	//	return nil, err
	//}
	//
	//if err := mixWallet(txM.Wallet, func(data dict.Dict) {
	//	doc.Mix("wallet", data)
	//}); err != nil {
	//	return nil, err
	//}
	//
	//if err := mixToken(txM.Mint, func(data dict.Dict) {
	//	doc.Mix("token", data)
	//}); err != nil {
	//	return nil, err
	//}

	return doc, nil
}
