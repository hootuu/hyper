package hyperplt

import (
	"github.com/hootuu/helix/components/hidem"
	"github.com/hootuu/helix/components/zplt"
)

func Idem() hidem.Factory {
	return zplt.UniIdem()
}
