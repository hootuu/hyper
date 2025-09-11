package supplyord

import "github.com/hootuu/hyper/hiorder"

const (
	Paying    = hiorder.Consensus
	Shipping  = hiorder.Executing
	Completed = hiorder.Completed
)
