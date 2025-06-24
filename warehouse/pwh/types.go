package pwh

type ID = uint64

type Direction byte

const (
	DirectionIn  Direction = 1
	DirectionOut Direction = 0

	DirectionLock   Direction = 11
	DirectionUnlock Direction = 10
)
