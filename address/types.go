package address

import (
	"github.com/hootuu/helix/components/sattva"
	"github.com/hootuu/hyper/address/maps"
)

type Location struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type Contact struct {
	Name string `json:"name"`
	Mobi string `json:"mobi"`
}

type Region struct {
	ID       RegionID  `json:"id"`
	Map      maps.Map  `json:"map"`
	Code     string    `json:"code"`
	Name     string    `json:"name"`
	Children []*Region `json:"children"`
}

type Address struct {
	ID       string                `json:"id"`
	Owner    sattva.Identification `json:"owners"`
	Region   RegionID              `json:"region"`
	Address  string                `json:"address"`
	FullAddr string                `json:"full_addr"`
	Contact  Contact               `json:"contact"`
	Default  bool                  `json:"default"`
	Location Location              `json:"location"`
	Province string                `json:"province"`
	City     string                `json:"city"`
	District string                `json:"district"`
	Tag      string                `json:"tag"`
}
