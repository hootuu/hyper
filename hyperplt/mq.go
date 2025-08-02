package hyperplt

import (
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/helix/unicom/hmq/hmq"
)

func MQ() *hmq.MQ {
	return gMQ
}

var gMQ *hmq.MQ

func init() {
	helix.AfterStartup(func() {
		gMQ = zplt.HelixMQ()
	})
}

func MqPublish(topic hmq.Topic, payload hmq.Payload) error {
	producer, err := zplt.HelixMqProducer()
	if err != nil {
		return err
	}
	return producer.Publish(topic, payload)
}
