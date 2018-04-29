/*
@Time : 18-4-20 下午4:55 
@Author : wangruixin
@File : contruct.go
*/

package construct

import (
	"time"

	"github.com/wrxcode/deploy-server/docker"
	"github.com/wrxcode/deploy-server/commit"
	"github.com/wrxcode/deploy-server/models"

	log "github.com/sirupsen/logrus"
	"fmt"
	"github.com/wrxcode/deploy-server/common/g"
)

const (
	CONTRUCT_SUCC = 2
	CONTRUCE_FAIL = 3
)

const (
	ContructType 	= 	0
)

const (
	GIT_LOGERROR = "Git repository address error, please check the input address!"
)

//构建镜像的入口函数
func ContructImage(dataId int64) {
	// 通过Nsq发送过来的构建记录表ID，来获取构建记录信息
	record, recordErr := models.ConstructRecord{}.GetById(dataId)
	if recordErr != nil {
		log.Errorf("read sql error: OrderType[%d] , DataId[%d], ErrorReason[%s]", ContructType, dataId, recordErr.Error())
		return
	}
	if record == nil{
		log.Errorf("read sql error: OrderType[%d] , DataId[%d], ErrorReason[no id in sql]", ContructType, dataId)
		return
	}

	// 通过构建记录表的工程ID，来获取工程信息
	project, projectErr := models.Project{}.GetById(record.ProjectId)
	if projectErr != nil {
		log.Errorf("read sql error: OrderType[%d] , DataId[%d], ErrorReason[%s]", ContructType, dataId, projectErr.Error())
		return
	}
	if project == nil{
		log.Errorf("read sql error: OrderType[%d] , DataId[%d], ErrorReason[no id in sql]", ContructType, dataId)
		return
	}

	// 获取提交的git版本
	commitKey, describe, commitErr := commit.GetCommit(project.GitDockerPath)
	if commitErr != nil {
		log.Errorf("Git repository address is error: ret[%s], reason[%s]", commitKey, commitErr)
		writeDatabaseBack(record, CONTRUCE_FAIL, 0, GIT_LOGERROR)
		return
	}

	//通过docker file + docker API 进行部署
	err := docker.DockerBuild(project.GitDockerPath, project.ProjectName, commitKey, dataId)
	if err != nil {
		log.Errorf("docker build err: %v; id[%d], commit[%s]", err, project.Id, commitKey)
		writeDatabaseBack(record, CONTRUCE_FAIL, 0, "")
		return
	}

	str, err := docker.DockerPush(project.ProjectName, commitKey)
	if err != nil {
		log.Fatalf("docker push err: ret[%s], reason[%v]", str, err)
		writeDatabaseBack(record, CONTRUCE_FAIL, 0, str)
		return
	}

	//镜像名
	var repo string
	if g.Conf().Repo.IsHost == 1 {
		repo = g.Conf().Repo.Host
	} else if g.Conf().Repo.IsIp == 1 {
		repo = g.Conf().Repo.Ip
	}
	if repo == "" {
		log.Fatalf("config file error:[repo]")
		return
	}
	repo = fmt.Sprintf("%s:%s", repo, g.Conf().Repo.Port)
	//拼接镜像名与私有仓库名，方便docker push使用
	repoData := fmt.Sprintf("%s/%s", repo, project.ProjectName)

	mirror := models.Mirror{0, repoData,commitKey, describe}
	mirrorId, err := models.Mirror{}.Add(&mirror)
	if err != nil {
		log.Errorf("add sql error: OrderType[%d] , DataId[%d], ErrorReason[%v]", ContructType, dataId, err)
		return
	}

	writeDatabaseBack(record, CONTRUCT_SUCC, mirrorId, "")
	return
}

// 回写数据库
func writeDatabaseBack(construct *models.ConstructRecord, status int, mirrorId int64, constructLog string) {
	if status == CONTRUCT_SUCC {
		construct.MirrorId = mirrorId
	} else {
		construct, recordErr := models.ConstructRecord{}.GetById(construct.Id)
		if recordErr != nil {
			log.Errorf("read sql error: OrderType[%d] , DataId[%d], ErrorReason[%s]", ContructType, construct.Id, recordErr.Error())
			return
		}
		if construct == nil{
			log.Errorf("read sql error: OrderType[%d] , DataId[%d], ErrorReason[no id in sql]", ContructType, construct.Id)
			return
		}
		construct.ConstructLog += constructLog
	}

	construct.ConstructStatu = status
	construct.ConstructEnd 	 = time.Now().Unix()

	err := models.ConstructRecord{}.Update(construct)
	if err != nil {
		log.Errorf("Database write back error: %s", err.Error())
	}
}