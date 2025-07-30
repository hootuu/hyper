package hiorder

import "github.com/hootuu/hyle/hypes/ex"

type Matter = interface {
	GetDigest() ex.Meta
}
