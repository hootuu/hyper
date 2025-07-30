package hyperplt

import (
	"github.com/hootuu/helix/components/zplt/zticktock"
	"github.com/hootuu/helix/ticktock"
)

func Ticktock() *ticktock.Worker {
	return zticktock.Ticktock()
}

func Postman() *ticktock.Postman {
	return zticktock.Postman()
}
