package hitopup

import (
	"github.com/hootuu/hyle/hypes/ex"
	"github.com/nineora/harmonic/chain"
)

type Matter struct {
	InAccount chain.Address `json:"in_account"`
}

func (m Matter) GetDigest() ex.Meta {
	return ex.Meta{
		"in_account": m.InAccount,
	}
}
