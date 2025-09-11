package prodord

import "github.com/hootuu/hyper/hiorder"

const (
	Initial  = hiorder.Initial
	Paying   = hiorder.Consensus
	Shipping = hiorder.Executing
	Timeout  = hiorder.Timeout
)
