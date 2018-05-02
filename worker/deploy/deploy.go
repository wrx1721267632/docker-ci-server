/*
@Time : 18-4-20 下午4:57
@Author : wangruixin
@File : deploy.go
*/

package deploy

import (
	log "github.com/sirupsen/logrus"
	"github.com/wrxcode/deploy-server/models"
)

const (
	ContructType = 1
)

func Deploy(dataId int64) {
	//通过dataID获取到对应部署日志表的信息
	deploy, deployErr := models.Deploy{}.GetById(dataId)
	if deployErr != nil {
		log.Errorf("read deploy record sql error: OrderType[%d] , DataId[%d], ErrorReason[%s]", ContructType, dataId, deployErr.Error())
		return
	}
	if deploy == nil {
		log.Errorf("read deploy record sql error: OrderType[%d] , DataId[%d], ErrorReason[no id in sql]", ContructType, dataId)
		return
	}

}
