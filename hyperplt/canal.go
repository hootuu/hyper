package hyperplt

import (
	"github.com/hootuu/helix/components/zplt/zcanal"
	"github.com/hootuu/helix/storage/hcanal"
)

func Canal() *hcanal.Canal {
	return zcanal.HelixCanal()
}
