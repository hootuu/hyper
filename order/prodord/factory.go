package prodord

import (
	"github.com/hootuu/hyper/hiorder"
)

type Factory struct {
	core *hiorder.Factory[Matter]
}

func NewFactory(code hiorder.Code) *Factory {
	f := &Factory{}
	f.core = hiorder.NewFactory[Matter](newDealer(code, f))
	return f
}

func (f *Factory) Core() *hiorder.Factory[Matter] {
	return f.core
}
