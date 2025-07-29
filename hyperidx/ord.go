package hyperidx

import (
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/helix/storage/hmeili"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyper/hiorder"
	"github.com/hootuu/hyper/hpay/payment"
	"github.com/hootuu/hyper/hshipping/shipping"
	"github.com/hootuu/hyper/hyperplt"
	"github.com/meilisearch/meilisearch-go"
	"github.com/nineora/harmonic/nineidx"
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
		"payment_id",
		"shipping_id",
		"tag",
		"ctrl",
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
	if m.Payer != "" {
		doc["payer"] = idx.GetPayerDigest(m.Payer)
	}
	doc["payee"] = m.Payee.MustToDict()
	doc["amount"] = m.Amount
	doc["status"] = m.Status
	doc["matter"] = m.Matter
	doc["payment_id"] = m.PaymentID
	if m.PaymentID != 0 {
		doc["payment"], err = idx.GetPaymentDigest(m.PaymentID)
		if err != nil {
			return nil, err
		}
	}
	doc["shipping_id"] = m.ShippingID
	if m.ShippingID != 0 {
		doc["shipping"], err = idx.GetShippingDigest(m.ShippingID)
		if err != nil {
			return nil, err
		}
	}
	doc["tag"] = m.Tag
	doc["ctrl"] = m.Ctrl
	doc["meta"] = m.Meta
	//doc["ex_status"] = m.ex //todo

	return doc, nil
}

func (idx *TxOrdIndexer) GetPayerDigest(payer collar.Link) map[string]any {
	payerM := payer.MustToDict()
	return map[string]any{
		"code": payerM.Code,
		"id":   payerM.ID,
		"info": nineidx.FastGetSattva(payerM.ID),
	}
}

func (idx *TxOrdIndexer) GetShippingDigest(shippingID shipping.ID) (map[string]any, error) {
	shipM, err := hdb.MustGet[shipping.ShipM](hyperplt.DB(), "id = ?", shippingID)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"uni_link":     shipM.UniLink.MustToDict(),
		"address":      shipM.Address,
		"courier_code": shipM.CourierCode,
		"tracking_no":  shipM.TrackingNo,
		"shipped_at":   shipM.ShippedAt,
	}, nil
}

func (idx *TxOrdIndexer) GetPaymentDigest(paymentID payment.ID) (map[string]any, error) {
	payM, err := hdb.MustGet[payment.JobM](hyperplt.DB(), "payment_id = ? AND payment_seq = ?", paymentID, 1)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"status": payM.Status,
		"ctx":    payM.Ctx,
		"pay_no": payM.PayNo,
	}, nil
}
