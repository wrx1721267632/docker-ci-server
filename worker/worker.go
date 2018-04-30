/*
@Time : 18-4-20 下午4:07 
@Author : wangruixin
@File : worker.go
*/

package worker

import (
	//"encoding/json"
	//"os"
	//"path/filepath"
	//"strconv"

	//"github.com/wrxcode/deploy-server/common/components"
	//"github.com/wrxcode/deploy-server/common/g"
	//"github.com/wrxcode/deploy-server/models"
	"github.com/wrxcode/deploy-server/worker/construct"
	"github.com/wrxcode/deploy-server/worker/deploy"
	"github.com/wrxcode/deploy-server/worker/rollback"

	log "github.com/sirupsen/logrus"

)

type Worker struct {
	OrderType	int 	`json:"order_type"`
	DataId 	  	int64 	`json:"data_id"`
}

const (
	ContructType 	= 	0
	Deploy			= 	1
	RollBack		=	2
)

// 具体指令处理函数
func (this *Worker) DoWorker() {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("Worker Error : %v", err)
		}
	}()

	switch this.OrderType {
		//构建处理
	case ContructType:
		//fmt.Println("contruct!!!",this.DataId, "\n\n\n")
		construct.ContructImage(this.DataId)
		break

		//部署处理
	case Deploy:
		//fmt.Println("deploy!!!",this.DataId, "\n\n\n")
		deploy.Deploy()
		break

		//回滚处理
	case RollBack:
		//fmt.Println("rooback!!!",this.DataId, "\n\n\n")
		rollback.RollBack()
		break
	default:
		log.Errorf("OrderType error: OrderType[%d] , DataId[%d]", this.OrderType, this.DataId)
	}
}