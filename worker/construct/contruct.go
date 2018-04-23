/*
@Time : 18-4-20 下午4:55 
@Author : wangruixin
@File : contruct.go
*/

package construct

import (
	"github.com/wrxcode/deploy-server/docker"
	"github.com/wrxcode/deploy-server/worker"
	"github.com/wrxcode/deploy-server/commit"
	"github.com/hjhcode/deploy-web/models"

	log "github.com/sirupsen/logrus"
)

//构建镜像的入口函数
func ContructImage(dataId int64) {
	// 通过Nsq发送过来的构建记录表ID，来获取构建记录信息
	record, recordErr := models.ConstructRecord{}.GetById(dataId)
	if recordErr != nil {
		log.Errorf("read sql error: OrderType[%d] , DataId[%d], ErrorReason[%s]", worker.ContructType, dataId, recordErr.Error())
		return
	}
	if record == nil{
		log.Errorf("read sql error: OrderType[%d] , DataId[%d], ErrorReason[no id in sql]", worker.ContructType, dataId)
		return
	}

	// 通过构建记录表的工程ID，来获取工程信息
	project, projectErr := models.Project{}.GetById(record.ProjectId)
	if projectErr != nil {
		log.Errorf("read sql error: OrderType[%d] , DataId[%d], ErrorReason[%s]", worker.ContructType, dataId, projectErr.Error())
		return
	}
	if project == nil{
		log.Errorf("read sql error: OrderType[%d] , DataId[%d], ErrorReason[no id in sql]", worker.ContructType, dataId)
		return
	}

	// 获取提交的git版本
	commitKey, commitErr := commit.GetCommit(project.GitDockerPath)

	//通过docker file + docker API 进行部署
	docker.DockerBuild()
}