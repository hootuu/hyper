package giftord

import (
	"context"
	"errors"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyle/hypes/ex"
	"github.com/hootuu/hyper/hiorder"
	"github.com/hootuu/hyper/hyperplt"
	"go.uber.org/zap"
	"time"
)

type CreateParas struct {
	Matter
	Idem  string      `json:"idem"`
	Payer collar.Link `json:"payer"`
	Payee collar.Link `json:"payee"`
	Title string      `json:"title"`
	Ex    *ex.Ex      `json:"ex"`
}

func (paras *CreateParas) Validate() error {
	if paras.Idem == "" {
		return errors.New("idem is empty")
	}
	if paras.Payer == "" {
		return errors.New("player is empty")
	}
	if paras.Title == "" {
		return errors.New("title is empty")
	}
	if len(paras.Items) == 0 {
		return errors.New("items is empty")
	}
	if paras.Amount == 0 {
		return errors.New("amount is empty")
	}
	for _, item := range paras.Items {
		if err := item.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (f *Factory) CreateAndPay(ctx context.Context, paras *CreateParas) (*hiorder.Order[Matter], error) {
	if err := paras.Validate(); err != nil {
		return nil, err
	}
	hlog.Info("giftord create order", zap.Any("paras", paras))
	engine, err := f.core.New(ctx, &hiorder.CreateParas[Matter]{
		Idem:    paras.Idem,
		Title:   paras.Title,
		Payer:   paras.Payer,
		Payee:   paras.Payee,
		Amount:  0,
		Payment: nil,
		Link:    collar.Build(f.core.Code(), paras.Payer.MustToID()).Link(),
		Matter: Matter{
			Items:  paras.Items,
			Amount: paras.Amount,
			Count:  paras.Count,
		},
		Ex: paras.Ex,
	})
	if err != nil {
		return nil, err
	}
	err = engine.Submit(ctx)
	if err != nil {
		return nil, err
	}

	err = hdb.Update[hiorder.OrderM](hyperplt.Tx(ctx), map[string]any{
		"status":         hiorder.Consensus,
		"consensus_time": time.Now(),
	}, "id = ?", engine.GetOrder().ID)
	if err != nil {
		return nil, err
	}

	order := engine.GetOrder()
	return order, nil
}
