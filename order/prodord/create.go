package prodord

import (
	"context"
	"errors"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyle/hypes/ex"
	"github.com/hootuu/hyper/hiorder"
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

func (f *Factory) Create(ctx context.Context, paras *CreateParas) (*hiorder.Order[Matter], error) {
	if err := paras.Validate(); err != nil {
		return nil, err
	}
	engine, err := f.core.New(ctx, &hiorder.CreateParas[Matter]{
		Idem:    paras.Idem,
		Title:   paras.Title,
		Payer:   paras.Payer,
		Payee:   paras.Payee,
		Amount:  paras.Amount,
		Payment: nil,
		Matter: Matter{
			Items:  paras.Items,
			Amount: paras.Amount,
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
	return engine.GetOrder(), nil
}
