package hyperidxsync

import (
	"github.com/hootuu/helix/components/zplt/zcanal"
	"github.com/hootuu/helix/storage/hcanal"
)

func canal() *hcanal.Canal {
	return zcanal.HelixCanal()
}
