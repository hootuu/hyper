package strategy

import (
	"fmt"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hmath"
)

type Pricing interface {
	GetType() string
	Validate() error
	ToBytes() []byte
}

const PricingProfit = "profit"

type ProfitPricing struct {
	Profit    hmath.Rate `json:"profit"`
	Precision int        `json:"precision"`
}

func (p *ProfitPricing) GetType() string {
	return PricingProfit
}

func (p *ProfitPricing) Validate() error {
	if p.Profit.Rate() >= 1 {
		return fmt.Errorf("profit rate must be less than 1")
	}
	if p.Precision <= 0 || p.Precision > 2 {
		return fmt.Errorf("precision must be between 0 and 2")
	}
	return nil
}

func (p *ProfitPricing) ToBytes() []byte {
	return hjson.MustToBytes(p)
}

func DefaultPricing() Pricing {
	return &ProfitPricing{
		Profit:    hmath.NewRate(hmath.Rate10000, 3000),
		Precision: 2,
	}
}
