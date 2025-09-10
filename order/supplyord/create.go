package supplyord

import (
	"context"
	"errors"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyle/hypes/ex"
	"github.com/hootuu/hyper/hiorder"
	"github.com/hootuu/hyper/hiprod/prod"
	"github.com/hootuu/hyper/hyperplt"
	"gorm.io/gorm"
)

type CreateParas struct {
	Matter
	Idem  string      `json:"idem"`
	Payer collar.Link `json:"payer"`
	Payee collar.Link `json:"payee"`
	Title string      `json:"title"`
	Link  collar.Link `json:"link"`
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
	return nil
}

func (f *Factory) Create(ctx context.Context, paras *CreateParas) (*hiorder.Order[Matter], error) {
	if err := paras.Validate(); err != nil {
		return nil, err
	}
	var spuIds []uint64
	for _, item := range paras.Items {
		if err := item.Validate(); err != nil {
			return nil, err
		}
		spuIds = append(spuIds, item.SpuID)
	}

	tx := hyperplt.DB()
	list, err := hdb.Find[prod.SpuM](func() *gorm.DB {
		return tx.Where("id in ?", spuIds)
	})
	if err != nil {
		return nil, err
	}
	spuPriceMap := make(map[prod.SpuID]uint64)
	for _, spu := range list {
		spuPriceMap[spu.ID] = spu.Cost
	}
	var amount uint64
	items := make([]*Item, 0)
	for _, item := range paras.Items {
		elems := &Item{
			SpuID:    item.SpuID,
			SkuID:    item.SkuID,
			VwhID:    item.VwhID,
			PwhID:    item.PwhID,
			Price:    spuPriceMap[item.SpuID],
			Quantity: item.Quantity,
			Amount:   spuPriceMap[item.SpuID] * item.Quantity,
		}
		items = append(items, elems)
		amount += elems.Amount
	}

	engine, err := f.core.New(ctx, &hiorder.CreateParas[Matter]{
		Idem:   paras.Idem,
		Title:  paras.Title,
		Payer:  paras.Payer,
		Payee:  paras.Payee,
		Amount: amount,
		//Payment: []payment.JobDefine{&payment.ThirdJob{
		//	ThirdCode: paras.PayChannelCode,
		//	Amount:    paras.Amount,
		//}},
		Link: paras.Link,
		Matter: Matter{
			Items:  items,
			Amount: amount,
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
	order := engine.GetOrder()

	return order, nil
}
