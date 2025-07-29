package prodord

import (
	"context"
	"github.com/hootuu/hyle/hcoin"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyle/hypes/ex"
	"github.com/hootuu/hyper/hiorder"
	"github.com/hootuu/hyper/hiprod/prod"
)

type CreateParas struct {
	Idem     string       `json:"idem"`
	Title    string       `json:"title"`
	ProdID   prod.ID      `json:"prod_id"`
	Payer    collar.Link  `json:"payer"`
	Payee    collar.Link  `json:"payee"`
	Quantity uint32       `json:"quantity"`
	Amount   hcoin.Amount `json:"amount"`
	Ex       *ex.Ex       `json:"ex"`
}

func (paras *CreateParas) Validate() error {
	//todo add validate
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
