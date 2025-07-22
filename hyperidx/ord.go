package hyperidx

import (
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/helix/storage/hmeili"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyper/hiorder"
	"github.com/hootuu/hyper/hyperplt"
	"github.com/meilisearch/meilisearch-go"
	"go.uber.org/zap"
)

const (
	OrdIndex      = "hyper_order"
	ordIdxVersion = "1_0_0"
)

type TxOrdIndexer struct{}

func (idx *TxOrdIndexer) GetName() string {
	return OrdIndex
}

func (idx *TxOrdIndexer) GetVersion() string {
	return ordIdxVersion
}

func (idx *TxOrdIndexer) Setting(index meilisearch.IndexManager) error {
	filterableAttributes := []string{
		"auto_id",
		"id",
		"code",
		"payer_code",
		"payer",
		"payer_acc_code",
		"payer_acc_id",
		"payee_code",
		"payee_id",
		"payee_acc_code",
		"payee_acc_id",
		"status",
		"timestamp",
		"title",
	}
	_, err := index.UpdateFilterableAttributes(&filterableAttributes)
	if err != nil {
		hlog.Err("hyper.idx.Setting: Error updating filterable attributes", zap.Error(err))
		return err
	}

	sortableAttributes := []string{
		"auto_id",
		"timestamp",
		"status",
	}
	_, err = index.UpdateSortableAttributes(&sortableAttributes)
	if err != nil {
		hlog.Err("hyper.idx.Setting: Error updating sortable attributes", zap.Error(err))
		return err
	}
	return nil
}

func (idx *TxOrdIndexer) Load(autoID int64) (hmeili.Document, error) {
	m, err := hdb.MustGet[hiorder.OrderM](hyperplt.DB(), "auto_id = ?", autoID)
	if err != nil {
		return nil, err
	}
	doc := hmeili.NewMapDocument(m.ID, m.AutoID, m.UpdatedAt.UnixMilli())
	doc["code"] = m.Code
	doc["title"] = m.Title
	doc["payer"] = m.Payer.MustToDict()
	doc["payee"] = m.Payee.MustToDict()
	doc["amount"] = m.Amount
	doc["status"] = m.Status
	doc["matter"] = m.Matter
	//todo

	return doc, nil
}
