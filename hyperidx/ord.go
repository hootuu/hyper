package hyperidx

import (
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/helix/storage/hmeili"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyper/hiorder"
	"github.com/hootuu/hyper/hiprod/pwh"
	"github.com/hootuu/hyper/hyperplt"
	"github.com/hootuu/hyper/payment"
	"github.com/hootuu/hyper/shipping"
	"github.com/meilisearch/meilisearch-go"
	"github.com/nineora/harmonic/nineidx"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

const (
	OrdIndex      = "hyper_order"
	ordIdxVersion = "1_0_4"
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
		"payee",
		"payee_acc_code",
		"payee_acc_id",
		"status",
		"timestamp",
		"title",
		"payment",
		"payment_id",
		"shipping_id",
		"consensus_ts",
		"tag",
		"ctrl",
		"created_at_ts",
		"supplier_id",
		"is_promotion",
		"category",
		"gjj_status",
		"user_order_id",
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

	searchableAttributes := []string{
		"title",
	}
	_, err = index.UpdateSettings(hmeili.BuildSearchSettings(searchableAttributes))
	if err != nil {
		hlog.Err("hyper.spu.Setting: Error updating settings", zap.Error(err))
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
	doc["created_at"] = m.CreatedAt
	doc["created_at_ts"] = m.CreatedAt.Unix()
	doc["code"] = m.Code
	doc["title"] = m.Title
	if m.Payer != "" {
		doc["payer"] = idx.GetPayerDigest(m.Payer)
	}
	//if m.Payee != "" {
	//	doc["payee"] = idx.GetPayerDigest(m.Payee)
	//}
	doc["amount"] = m.Amount
	doc["status"] = m.Status
	doc["matter"] = m.Matter
	doc["consensus_time"] = m.ConsensusTime
	if m.ConsensusTime != nil {
		doc["consensus_ts"] = m.ConsensusTime.Unix()
	}
	doc["executing_time"] = m.ExecutingTime
	doc["canceled_time"] = m.CanceledTime
	doc["completed_time"] = m.CompletedTime
	doc["timeout_time"] = m.TimeoutTime
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
	var supplierId string
	meta := *hjson.MustFromBytes[dict.Dict](m.Meta)
	if len(meta) > 0 {
		supplierId = meta.Get("product.supplier_id").String()
		if supplierId == "" {
			supplierId = meta.Get("game.forge.supplier_id").String()
		}
		doc["is_promotion"] = meta.Get("is_promotion").Bool()
		doc["category"] = cast.ToInt64(meta.Get("category").Data())

		obj := meta.Get("gjj_info").MSI()
		if len(obj) > 0 {
			doc["gjj_status"] = true
		} else {
			doc["gjj_status"] = false
		}

		// n19专用查询用户单号
		doc["user_order_id"] = meta.Get("orderId").String()
	}
	doc["supplier_id"] = supplierId

	doc["tag"] = m.Tag
	doc["ctrl"] = m.Ctrl
	doc["meta"] = m.Meta
	//doc["ex_status"] = m.ex //todo

	return doc, nil
}

func (idx *TxOrdIndexer) GetPayerDigest(payer collar.Link) map[string]any {
	payerM := payer.MustToDict()
	data := map[string]any{
		"code": payerM.Code,
		"id":   payerM.ID,
	}
	if payerM.Code == "SATTVA" {
		info := nineidx.FastGetSattva(payerM.ID)
		data["info"] = info
	} else if payerM.Code == "PWH" {
		pwhInfo, _ := pwh.MustGetById(cast.ToUint64(payerM.ID))
		if pwhInfo != nil {
			data["name"] = pwhInfo.Memo
		}
	}
	return data
}

func (idx *TxOrdIndexer) GetShippingDigest(shippingID shipping.ID) (map[string]any, error) {
	shipM, err := hdb.MustGet[shipping.ShipM](hyperplt.DB(), "id = ?", shippingID)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"address":        shipM.Address,
		"courier_code":   shipM.CourierCode,
		"tracking_no":    shipM.TrackingNo,
		"submitted_time": shipM.SubmittedTime,
		"timeout":        shipM.Timeout,
	}, nil
}

func (idx *TxOrdIndexer) GetPaymentDigest(paymentID payment.ID) (map[string]any, error) {
	jobM, err := hdb.MustGet[payment.JobM](hyperplt.DB(), "payment_id = ? AND payment_seq = ?", paymentID, 1)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"ctx":            jobM.Ctx,
		"pay_no":         jobM.PayNo,
		"prepared_time":  jobM.PreparedTime,
		"canceled_time":  jobM.CanceledTime,
		"timeout_time":   jobM.TimeoutTime,
		"completed_time": jobM.CompletedTime,
	}, nil
}
