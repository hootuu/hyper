package hiorder

import (
	"context"
	"errors"
	"fmt"
	"github.com/hootuu/helix/components/hnid"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/storage/hdb"
	"github.com/hootuu/hyle/data/ctrl"
	"github.com/hootuu/hyle/data/dict"
	"github.com/hootuu/hyle/data/tag"
	"github.com/hootuu/hyle/hcoin"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hypes/collar"
	"github.com/hootuu/hyle/hypes/ex"
	"github.com/hootuu/hyper/hpay"
	"github.com/hootuu/hyper/hpay/payment"
	"github.com/hootuu/hyper/hyperplt"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
)

type Factory[T Matter] struct {
	dealer      Dealer[T]
	idGenerator hnid.Generator
}

func NewFactory[T Matter](dealer Dealer[T]) *Factory[T] {
	f := &Factory[T]{dealer: dealer}
	helix.Use(f.helix())
	return f
}

type CreateParas[T Matter] struct {
	Title   string              `json:"title"`
	Payer   collar.Link         `json:"payer"`
	Payee   collar.Link         `json:"payee"`
	Amount  hcoin.Amount        `json:"amount"`
	Payment []payment.JobDefine `json:"payment"`
	Matter  T                   `json:"matter"`
	Ctrl    ctrl.Ctrl           `json:"ctrl"`
	Tag     tag.Tag             `json:"tag"`
	Meta    dict.Dict           `json:"meta"`
	UniLink collar.Link         `json:"uniLink"`
}

func (p *CreateParas[T]) Validate() error {
	return nil // todo TODOTODO
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

	ordID := f.nextID()
	ord := Order[T]{
		ID:        ordID,
		Code:      f.Code(),
		Title:     paras.Title,
		Payer:     paras.Payer,
		Payee:     paras.Payee,
		Matter:    paras.Matter,
		Amount:    paras.Amount,
		PaymentID: 0,
		Status:    Draft,
		Ctrl:      paras.Ctrl,
		Tag:       paras.Tag,
		Meta:      paras.Meta,
	}

	if len(paras.Payment) > 0 {
		paymentID, err := hpay.Create(ctx, payment.CreateParas{
			Payer:   paras.Payer,
			Payee:   paras.Payee,
			BizLink: f.OrderCollar(ord.ID).Link(),
			Amount:  paras.Amount,
			Ex:      nil,
			Jobs:    paras.Payment,
		})
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

func (f *Factory[T]) SetPayment(
	ctx context.Context,
	ordId ID,
	paymentArr []payment.JobDefine,
	exM *ex.Ex,
) (paymentID payment.ID, err error) {
	if hlog.IsElapseComponent() {
		defer hlog.ElapseWithCtx(ctx, "hyper.order.SetPayment",
			hlog.F(zap.Uint64("ordId", ordId)),
			func() []zap.Field {
				if err != nil {
					return []zap.Field{zap.Error(err), zap.Any("payment", paymentArr)}
				}
				return nil
			},
		)()
	}
	if len(paymentArr) == 0 {
		return 0, fmt.Errorf("require payment")
	}
	ordM, err := hdb.MustGet[OrderM](hyperplt.Tx(ctx), "id = ?", ordId)
	if err != nil {
		hlog.Err("MustGet[OrderM] err", hlog.TraceInfo(ctx), zap.Error(err))
		return 0, errors.New("Load Order Failed: " + err.Error())
	}
	err = hdb.Tx(hyperplt.Tx(ctx), func(tx *gorm.DB) error {
		innerCtx := hdb.TxCtx(tx, ctx)
		paymentID, err = hpay.Create(innerCtx, payment.CreateParas{
			Payer:   ordM.Payer,
			Payee:   ordM.Payee,
			BizLink: f.OrderCollar(ordM.ID).Link(),
			Amount:  ordM.Amount,
			Ex:      exM,
			Jobs:    paymentArr,
		})
		err = payment.Prepare(innerCtx, paymentID)
		if err != nil {
			hlog.Err("payment.Prepare err", hlog.TraceInfo(ctx), zap.Error(err))
			return errors.New("payment prepared: " + err.Error())
		}
		ordMut := map[string]any{
			"payment_id": paymentID,
		}
		err = hdb.Update[OrderM](tx, ordMut, "id = ?", ordId)
		if err != nil {
			hlog.Err("Update[OrderM] err", hlog.TraceInfo(ctx), zap.Error(err))
			return errors.New("Update Order Failed: " + err.Error())
		}
		return nil
	})
	if err != nil {
		hlog.Err("Tx err", hlog.TraceInfo(ctx), zap.Error(err))
		return 0, errors.New("Tx Failed: " + err.Error())
	}
	return paymentID, nil
}

func (f *Factory[T]) Code() Code {
	return f.dealer.Code()
}

func (f *Factory[T]) OrderCollar(id ID) collar.Collar {
	return collar.Build(f.Code(), fmt.Sprintf("%d", id))
}

func (f *Factory[T]) onPaymentAltered(payload *payment.PaymentPayload) (err error) {
	if payload == nil {
		hlog.Fix("hyper.f.onPaymentAlter: payload is nil")
		return nil
	}
	defer hlog.Elapse("hiorder.f.onPaymentAltered",
		hlog.F(zap.String("ord.id", payload.BizID)),
		func() []zap.Field {
			if err != nil {
				return []zap.Field{zap.Error(err)}
			}
			return nil
		})()
	ordID := cast.ToUint64(payload.BizID)
	ordM, err := hdb.Get[OrderM](hyperplt.DB(), "id = ?", ordID)
	if err != nil {
		return err
	}
	if ordM == nil {
		hlog.Fix("hyper.f.onPaymentAlter: order not found", zap.Uint64("id", ordID))
		return nil
	}

	err = f.dealer.OnPaymentAltered(&PaymentAltered[T]{
		Order:     orderMto[T](ordM),
		PaymentID: payload.PaymentID,
		SrcStatus: payload.Src,
		DstStatus: payload.Dst,
	})
	if err != nil {
		return err
	}

	return nil
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
			SetAutoInc(5, 1, 99999, 1000), //todo
	)
	if err != nil {
		return nil, err
	}
	err = hyperplt.DB().AutoMigrate(&OrderM{})
	if err != nil {
		return nil, err
	}
	fmt.Println("---------------REGISTER FMQ HANDLE", f.Code()) //todo
	//doRegFactoryMqHandle(fmt.Sprintf("HIORD_%s", f.Code()), f.onPaymentAltered)
	payment.MqRegisterPaymentAlter(f.Code(), func(ctx context.Context, payload *payment.PaymentPayload) error {
		return f.onPaymentAltered(payload)
	})
	return nil, nil
}

func (f *Factory[T]) doShutdown(_ context.Context) {
}
