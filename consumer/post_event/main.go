package main

import (
	"encoding/json"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/g10guang/graduation/constdef"
	"github.com/g10guang/graduation/consumer/post_event/handler"
	"github.com/g10guang/graduation/model"
	"github.com/g10guang/graduation/tools"
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
)

func main() {
	tools.InitLog()
	defer clean()

	config := nsq.NewConfig()
	config.MaxBackoffDuration = 0
	consumer, err := nsq.NewConsumer(constdef.PostFileEventTopic, "compress", config)
	if err != nil {
		panic(err)
	}
	consumer.ChangeMaxInFlight(200)
	consumer.AddConcurrentHandlers(nsq.HandlerFunc(compress), 10)

	if err = consumer.ConnectToNSQLookupds([]string{constdef.NsqLookupdAddr}); err != nil {
		logrus.Panicf("ConnectToNSQLookupds Error: %s", err)
		panic(err)
	}

	//logrus.Infof("ConnectToNsda: %v", []string{constdef.NsqdAddr})

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

// 将图片转化为 jpeg/png 格式
func compress(message *nsq.Message) error {
	logrus.Debugf("compress message: %+v", message)
	ctx := tools.NewCtxWithLogID()
	msg := parsePostFileEventMsg(message.Body)
	if msg == nil {
		return errors.New("message error")
	}
	h := handler.NewCompressHandler(msg)
	if err := h.Handle(ctx); err != nil {
		logrus.Errorf("CompressHandler Error: %s", err)
		return err
	}
	logrus.Infof("CompressHandler Success")
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
