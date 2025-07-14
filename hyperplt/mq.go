package hyperplt

import (
	"github.com/hootuu/helix/unicom/hmq/hmq"
	"github.com/hootuu/helix/unicom/hmq/hnsq"
)

func MQ() *hmq.MQ {
	return gMQ
}

var gMQ *hmq.MQ

func init() {
	gMQ = hmq.NewMQ("hyper_mq", hnsq.NewNsqMQ())
}
