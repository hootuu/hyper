package maps

import "github.com/hootuu/hyle/data/pagination"

type Map int

const (
	AMAP Map = 1
)

type Level int

const (
	Country  Level = 1
	Province Level = 2
	City     Level = 3
	District Level = 4
	Street   Level = 5
)

type Region struct {
	Map     Map    `json:"map"`
	Code    string `json:"code"`
	Name    string `json:"name"`
	Level   Level  `json:"level"`
	Address string `json:"address"`
}

type MapProvider interface {
	Region(page pagination.Page) (pagination.Pagination[Region], error)
}
