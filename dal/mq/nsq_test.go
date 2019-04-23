package mq

import (
	"github.com/g10guang/graduation/constdef"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestPublishNsq(t *testing.T) {
	err := nsqProducer.Ping()
	if err != nil {
		logrus.Errorf("NsqProducer Ping Error: %s", err)
	} else {
		logrus.Infof("NsqProducer Ping Success")
	}

	err = PublishNsq(constdef.PostFileEventTopic, []byte("hello world"))
	if err != nil {
		logrus.Errorf("PublishNsq Error: %s", err)
	} else {
		logrus.Infof("PublishNsq Success")
	}

}
