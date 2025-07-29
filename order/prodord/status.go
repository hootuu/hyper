package prodord

import "github.com/hootuu/hyper/hiorder"

const (
	_ hiorder.ExStatus = iota
	Paying
	Shipping
)
