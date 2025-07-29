package hiorder

import "github.com/hootuu/helix/helix"

func NewFactory[T Matter](dealer Dealer[T]) *Factory[T] {
	InitIfNeeded()
	f := buildFactory[T](dealer)
	helix.Use(f.helix())
	return f
}
