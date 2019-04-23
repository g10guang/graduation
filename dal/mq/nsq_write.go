package mq

import (
	"github.com/g10guang/graduation/constdef"
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
)

// when a file posted, send a message to nsq
var nsqProducer *nsq.Producer

func init() {
	var err error
	config := nsq.NewConfig()
	nsqProducer, err = nsq.NewProducer(constdef.NsqdAddr, config)
	if err != nil {
		panic(err)
	}
}

func StopNsqProducer() {
	nsqProducer.Stop()
}

func PublishNsq(topic string, content []byte) error {
	logrus.Debugf("PublishNsq topic: %s content: %s", topic, string(content))
	err := nsqProducer.Publish(topic, content)
	if err != nil {
		logrus.Errorf("PublishNsq Error: %s topic: %s content: %s", err, topic, string(content))
	}
	return err
}
