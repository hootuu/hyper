package hiorder

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/components/hnid"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyper/hyperplt"
	"github.com/hootuu/hyper/payment"
	"github.com/hootuu/hyper/shipping"
	"go.uber.org/zap"
	"strings"
)

type Factory[T Matter] struct {
	dealer          Dealer[T]
	idGenerator     hnid.Generator
	alterHandlerArr []func(ctx context.Context, payload *AlterPayload) error
}

func buildFactory[T Matter](dealer Dealer[T]) *Factory[T] {
	f := &Factory[T]{dealer: dealer}
	doInjectUniOrdAlterHandler(f.Code(), f.onOrdAlter)
	payment.ListeningAlter(f.Code(), f.onPaymentAlter)
	shipping.ListeningAlter(f.Code(), f.onShippingAlter)
	return f
}

func (f *Factory[T]) Load(ctx context.Context, id ID) (*Engine[T], error) {
	tx := hyperplt.Tx(ctx)
	ordM, err := hdb.Get[OrderM](tx, "id = ?", id)
	if err != nil {
		return nil, err
	}
	if ordM == nil {
		return nil, fmt.Errorf("order not found: %d", id)
	}
	ord := orderMto[T](ordM)
	dealer, err := f.dealer.Build(*ord)
	if err != nil {
		hlog.Err("hyper.order.Load: dealer.Builder", zap.Error(err))
		return nil, err
	}

	return newEngine(dealer, ord, f), nil
}

func (f *Factory[T]) Code() Code {
	return f.dealer.Code()
}

func (f *Factory[T]) OrderCollar(id ID) collar.Collar {
	return collar.Build(strings.ToUpper(f.Code()), fmt.Sprintf("%d", id))
}

func (f *Factory[T]) nextID() ID {
	return f.idGenerator.NextUint64()
}

func (f *Factory[T]) helix() helix.Helix {
	hCode := fmt.Sprintf("hyper_ord_%s", strings.ToLower(f.Code()))
	return helix.BuildHelix(hCode, f.doStartup, f.doShutdown)
}

func (f *Factory[T]) doStartup() (context.Context, error) {
	var err error
	iderCode := fmt.Sprintf("hyper_ord_%s_ider", strings.ToLower(f.Code()))
	f.idGenerator, err = hnid.NewGenerator(iderCode,
		hnid.NewOptions(1, 8).
			SetTimestamp(hnid.Minute, false).
			SetAutoInc(5, 1, 99999, 1000),
	)
	if err != nil {
		return nil, err
	}
	err = hyperplt.DB().AutoMigrate(&OrderM{})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (f *Factory[T]) doShutdown(_ context.Context) {
}
