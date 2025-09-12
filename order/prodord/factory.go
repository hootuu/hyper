package prodord

import (
	"github.com/hootuu/hyper/hiorder"
)

const Code = "PROD_ORDER"

type Factory struct {
	core *hiorder.Factory[Matter]
}

func newFactory() *Factory {
	f := &Factory{}
	f.core = hiorder.NewFactory[Matter](newDealer(Code, f))
	return f
}

func (f *Factory) Core() *hiorder.Factory[Matter] {
	return f.core
}

var factory *Factory

//var once sync.Once

func GetFactory() *Factory {
	if factory != nil {
		return factory
	}
	//once.Do(func() {
	factory = newFactory()
	//})
	return factory
}
