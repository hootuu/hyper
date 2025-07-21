package hyperplt

import (
	"github.com/hootuu/helix/components/zplt"
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

func MqPublish(topic hmq.Topic, payload hmq.Payload) error {
	producer, err := zplt.HelixMqProducer()
	if err != nil {
		return err
	}
	return producer.Publish(topic, payload)
}
