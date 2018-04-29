/*
@Time : 18-4-20 下午2:51
@Author : wangruixin
@File : consumer
*/

package dispatch

import (
	"github.com/wrxcode/deploy-server/common/g"
	"github.com/nsqio/go-nsq"

	log "github.com/sirupsen/logrus"
)

// Nsq 对应配置
type Consumer struct {
	NsqConsumer *nsq.Consumer
	Cfg         *nsq.Config
	Topic       string
	Channel     string
}

var consumers []*Consumer
var handlerCount chan int

// 启动Nsq
func StartConsume() {
	cfg := g.Conf()

	consumer := new(Consumer)
	go consumer.newConsumer(cfg.Nsq.DeployTopic, cfg.Nsq.DeployChannel)
	consumers = append(consumers, consumer)

	handlerCount = make(chan int, cfg.Nsq.HandlerCount)
}

// 停止Nsq
func StopConsume() {
	for _, c := range consumers {
		c.NsqConsumer.Stop()
	}
}

// 创建新的Nsq消费者
func (this *Consumer) newConsumer(topic, channel string) {
	this.Cfg = nsq.NewConfig()
	this.Topic = topic
	this.Channel = channel

	var err error
	this.NsqConsumer, err = nsq.NewConsumer(topic, channel, this.Cfg)
	if err != nil {
		log.Fatal(err)
	}

	this.NsqConsumer.AddHandler(&Handler{Topic: topic})

	err = this.NsqConsumer.ConnectToNSQLookupds(g.Conf().Nsq.Lookupds)
	if err != nil {
		log.Fatal(err)
	}
}

