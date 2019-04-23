package main

import (
	"encoding/json"
	"fmt"
	"github.com/g10guang/graduation/constdef"
	"github.com/g10guang/graduation/consumer/delete_event/handler"
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
	consumer, err := nsq.NewConsumer(constdef.DeleteFileEventTopic, "delete", config)
	consumer.ChangeMaxInFlight(200)
	consumer.AddHandler(nsq.HandlerFunc(delete_))

	if err = consumer.ConnectToNSQLookupds([]string{constdef.NsqLookupdAddr}); err != nil {
		logrus.Panicf("ConnectToNSQLookupds Error: %s", err)
		panic(err)
	}

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

func delete_(message *nsq.Message) error {
	ctx := tools.NewCtxWithLogID()
	msg := parseDeleteFileEventMsg(message.Body)
	if msg.Fid == 0 || msg.Uid == 0 {
		return fmt.Errorf("invalid message: %s", string(message.Body))
	}
	h := handler.NewDeleteStorageHandler(msg)
	if err := h.Handle(ctx); err != nil {
		logrus.Errorf("Delete Storage Error: %s", err)
		return err
	}
	return nil
}

func parseDeleteFileEventMsg(body []byte) *model.DeleteFileEvent {
	msg := &model.DeleteFileEvent{}
	err := json.Unmarshal(body, msg)
	if err != nil {
		logrus.Errorf("unmarshal DeleteFileEvent message Error: %s", err)
	}
	return msg
}
