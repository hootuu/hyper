package hiorder

import (
	"context"
	"errors"
	"fmt"
	"github.com/hootuu/hyle/hcoin"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyle/hypes/ex"
	"github.com/hootuu/hyper/hyperplt"
	"github.com/hootuu/hyper/payment"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type CreateParas[T Matter] struct {
	Idem    string              `json:"idem"`
	Title   string              `json:"title"`
	Payer   collar.Link         `json:"payer"`
	Payee   collar.Link         `json:"payee"`
	Amount  hcoin.Amount        `json:"amount"`
	Payment []payment.JobDefine `json:"payment"`
	Matter  T                   `json:"matter"`
	Link    collar.Link         `json:"link"`
	Ex      *ex.Ex              `json:"ex"`
}

func (p *CreateParas[T]) Validate() error {
	if p.Idem == "" {
		return errors.New("idem is required")
	}
	if p.Title == "" {
		return errors.New("title is required")
	}
	if p.Payer == "" {
		return errors.New("payer is required")
	}
	if p.Amount == 0 {
		return errors.New("amount is required")
	}
	return nil
}

func (f *Factory[T]) New(ctx context.Context, paras *CreateParas[T]) (engine *Engine[T], err error) {
	if paras == nil {
		return nil, fmt.Errorf("create paras is nil")
	}
	if err := paras.Validate(); err != nil {
		return nil, err
	}
	if hlog.IsElapseComponent() {
		defer hlog.ElapseWithCtx(ctx, "hiorder.Create", hlog.F(), func() []zap.Field {
			if err != nil {
				return []zap.Field{
					zap.Any("paras", paras),
					zap.Error(err),
				}
			}
			return nil
		})()
	}
	if err := hyperplt.Idem().MustCheck(paras.Idem); err != nil {
		return nil, err
	}

	ordID := f.nextID()
	ord := Order[T]{
		ID:        ordID,
		Code:      f.Code(),
		Title:     paras.Title,
		Payer:     paras.Payer,
		Payee:     paras.Payee,
		Matter:    paras.Matter,
		Amount:    paras.Amount,
		Link:      paras.Link,
		PaymentID: 0,
		Status:    Draft,
		Ex:        paras.Ex,
	}

	if len(paras.Payment) > 0 {
		paymentID, err := payment.Create(ctx, &payment.CreateParas{
			Idem:    fmt.Sprintf("HYPER:ORD:PAY:%s:%d", f.Code(), ordID),
			Payer:   paras.Payer,
			Payee:   paras.Payee,
			BizCode: f.Code(),
			BizID:   cast.ToString(ordID),
			Amount:  paras.Amount,
			Ex:      paras.Ex,
			Jobs:    paras.Payment,
			Timeout: 0,
		})
		if err != nil {
			return nil, err
		}
		err = payment.Prepare(ctx, paymentID)
		if err != nil {
			return nil, err
		}
		ord.PaymentID = paymentID
	}

	dealer, err := f.dealer.Build(ord)
	if err != nil {
		hlog.Err("hyper.order.Create: dealer.Builder", zap.Error(err))
		return nil, err
	}

	return newEngine(dealer, &ord, f), nil
}
