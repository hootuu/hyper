package shipping

import (
	"context"
	"errors"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/ex"
	"github.com/hootuu/hyper/hyperplt"
	"go.uber.org/zap"
	"time"
)

type CreateParas struct {
	Idem    string        `json:"idem"`
	BizCode string        `json:"biz_code"`
	BizID   string        `json:"biz_id"`
	Address *Address      `json:"address"`
	Ex      *ex.Ex        `json:"ex"`
	Timeout time.Duration `json:"timeout"`
}

func (paras *CreateParas) Validate() error {
	if paras.Idem == "" {
		return errors.New("idem is required")
	}
	if paras.BizCode == "" {
		return errors.New("biz_code is required")
	}
	if paras.BizID == "" {
		return errors.New("biz_id is required")
	}
	if paras.Address == nil {
		return errors.New("address is required")
	}
	if err := paras.Address.Validate(); err != nil {
		return err
	}
	if paras.Timeout == 0 {
		paras.Timeout = 7 * 24 * time.Hour
	}
	return nil
}

func Create(ctx context.Context, paras *CreateParas) (id ID, err error) {
	InitIfNeeded()
	if paras == nil {
		return 0, errors.New("paras is required")
	}
	if err := paras.Validate(); err != nil {
		return 0, err
	}
	if hlog.IsElapseComponent() {
		defer hlog.ElapseWithCtx(ctx, "hyper.shipping.Create",
			hlog.F(zap.String("bizCode", paras.BizCode)),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Any("paras", paras), zap.Error(err)}
				}
				return nil
			},
		)()
	}
	if err := hyperplt.Idem().MustCheck(paras.Idem); err != nil {
		return 0, err
	}

	id = nxtShippingID()
	exM := ex.MustEx(paras.Ex)
	shippingM := &ShipM{
		Template: hdb.Template{
			Ctrl: exM.Ctrl,
			Tag:  hjson.MustToBytes(exM.Tag),
			Meta: hjson.MustToBytes(exM.Meta),
		},
		ID:               id,
		BizCode:          paras.BizCode,
		BizID:            paras.BizID,
		Status:           Initialized,
		Address:          hjson.MustToBytes(paras.Address),
		Timeout:          paras.Timeout,
		TimeoutCompleted: false,
	}
	err = hdb.Create[ShipM](hyperplt.Tx(ctx), shippingM)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func Update(ctx context.Context, orderId, courierCode, trackingNo string) error {
	if orderId == "" {
		return errors.New("order_id is required")
	}
	if courierCode == "" {
		return errors.New("courier_code is required")
	}
	if trackingNo == "" {
		return errors.New("tracking_no is required")
	}
	return hdb.Update[ShipM](hyperplt.Tx(ctx), map[string]any{
		"courier_code": courierCode,
		"tracking_no":  trackingNo,
	}, "biz_id = ?", orderId)
}
