package dispatch

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/nsqio/go-nsq"
)

func TestSendMessCpp(t *testing.T) {
	Nsq{}.send("deploy", &SendMess{0,1})
}

type Nsq struct{}

type SendMess struct {
	OrderType	int 	`json:"order_type"`		//命令类型
	DataId 	  	int64 	`json:"data_id"`		//数据库ID
}

func (this Nsq) send(topic string, sendMess *SendMess) {
	if topic != "deploy" {
		err := errors.New("topic is false!")
		panic(err.Error())
	}
	adds := [1]string{"xupt2.fightcoder.com:9002"}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := r.Intn(len(adds))
	msg, err := json.Marshal(sendMess)
	if err != nil {
		fmt.Println(err)
	}
	postNsq(adds[num], topic, msg)
}
func postNsq(address, topic string, msg []byte) {
	config := nsq.NewConfig()
	if w, err := nsq.NewProducer(address, config); err != nil {
		panic(err)
	} else {
		w.Publish(topic, msg)
	}
}

