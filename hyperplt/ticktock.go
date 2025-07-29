package hyperplt

import (
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/ticktock"
)

func Ticktock() *ticktock.Worker {
	return zplt.Ticktock()
}

func Postman() *ticktock.Postman {
	return zplt.Postman()
}
