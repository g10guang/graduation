package main

import (
	"encoding/json"
	"errors"
	"github.com/g10guang/graduation/constdef"
	"github.com/g10guang/graduation/consumer/checksum/handler"
	"github.com/g10guang/graduation/model"
	"github.com/g10guang/graduation/tools"
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	tools.InitLog()
	defer clean()

	config := nsq.NewConfig()
	config.MaxBackoffDuration = 0
	consumer, err := nsq.NewConsumer(constdef.PostFileEventTopic, "checksum", config)
	if err != nil {
		panic(err)
	}
	consumer.ChangeMaxInFlight(200)
	consumer.AddConcurrentHandlers(nsq.HandlerFunc(checksum), 10)
	if err = consumer.ConnectToNSQLookupds([]string{constdef.NsqLookupdAddr}); err != nil {
		logrus.Panicf("ConnectToNSQLookupds Error: %s", err)
		panic(err)
	}
	//if err = consumer.ConnectToNSQDs([]string{constdef.NsqdAddr}); err != nil {
	//	logrus.Panicf("ConnectToNSQDs Error: %s", err)
	//	panic(err)
	//}

	shutdown := make(chan os.Signal, 2)
	signal.Notify(shutdown, syscall.SIGINT)

	for {
		select {
		case <-consumer.StopChan:
			goto exit
		case <-shutdown:
			consumer.Stop()
			goto exit
		}
	}

exit:
	logrus.Infof("consumer exit")
}

func clean() {

}

// 计算 md5 checksum
func checksum(message *nsq.Message) error {
	logrus.Debugf("checksum message: %+v", message)
	ctx := tools.NewCtxWithLogID()
	msg := parsePostFileEventMsg(message.Body)
	if msg == nil {
		return errors.New("message error")
	}
	h := handler.NewChecksumHandler(msg)
	if err := h.Handle(ctx); err != nil {
		logrus.Errorf("ChecksumHandler Error: %s", err)
		return err
	}
	logrus.Infof("ChecksumHandler Success")
	return nil
}

func parsePostFileEventMsg(body []byte) *model.PostFileEvent {
	logrus.Infof("post_file event message: %s", string(body))
	m := new(model.PostFileEvent)
	if err := json.Unmarshal(body, m); err != nil {
		logrus.Errorf("PostFileEvent message Error: %s", err)
		return nil
	}
	return m
}
