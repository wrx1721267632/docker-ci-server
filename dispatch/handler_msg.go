/*
@Time : 18-4-20 下午3:50 
@Author : wangruixin
@File : handler_msg.go
*/

package dispatch

import (
	"encoding/json"

	"github.com/wrxcode/deploy-server/worker"

	"github.com/nsqio/go-nsq"

	log "github.com/sirupsen/logrus"
)

type Handler struct {
	Topic string
}

func (this *Handler) HandleMessage(m *nsq.Message) error {
	log.Infof("HandbleMessage: ", string(m.Body))

	workerData := new(worker.Worker)
	if err := json.Unmarshal(m.Body, workerData); err != nil {
		log.Errorf("unmarshal JudgerData from NsqMessage failed, err: %v, event:%s", err, m.Body)
		return nil
	}

	log.Infof("consume Message from dispatch: %#v", workerData)

	handlerCount <- 1
	go this.doWorker(workerData)

	return nil
}

func (this *Handler) doWorker(workerData *worker.Worker) {
	defer func() {
		<-handlerCount
	}()

	workerData.DoWorker()
}

