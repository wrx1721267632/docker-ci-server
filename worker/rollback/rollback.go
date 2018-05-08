/*
@Time : 18-4-20 下午6:25
@Author : wangruixin
@File : rollback.go
*/

package rollback

import (
	"encoding/json"
	"fmt"

	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/wrxcode/deploy-server/docker"
	"github.com/wrxcode/deploy-server/models"
	"github.com/wrxcode/deploy-server/worker/deploy"
)

const (
	RollbackType = 2
)

func Rollback(dataId int64) {
	//通过dataID获取到对应部署日志表的信息
	deployInfo, deployErr := models.Deploy{}.GetById(dataId)
	if deployErr != nil {
		log.Errorf("read deploy record sql error: OrderType[%d] , DataId[%d], ErrorReason[%s]\n", RollbackType, dataId, deployErr.Error())
		return
	}
	if deployInfo == nil {
		log.Errorf("read deploy record sql error: OrderType[%d] , DataId[%d], ErrorReason[no id in sql]\n", RollbackType, dataId)
		return
	}

	//通过deploy中的serviceId找到对应service信息
	serviceInfo, serviceErr := models.Service{}.GetById(deployInfo.ServiceId)
	if serviceErr != nil {
		log.Errorf("read service sql error: OrderType[%d] , DataId[%d], ErrorReason[%s]\n", RollbackType, dataId, serviceErr.Error())
		return
	}
	if serviceInfo == nil {
		log.Errorf("read service sql error: OrderType[%d] , DataId[%d], ErrorReason[no id in sql]\n", RollbackType, dataId)
		return
	}

	//通过serviceId找到上一次部署成功的信息，用来作为回滚任务的快照
	lastSuccDeployInfo, lastSuccErr := models.Deploy{}.GetDeployBackData(serviceInfo.Id)
	if lastSuccErr != nil {
		log.Errorf("read service sql error: OrderType[%d] , DataId[%d], ErrorReason[%s]\n", RollbackType, dataId, serviceErr.Error())
		return
	}
	if lastSuccDeployInfo == nil {
		log.Errorf("read service sql error: OrderType[%d] , DataId[%d], ErrorReason[no id in sql]\n", RollbackType, dataId)
		deploy.RewriteDeployLog(dataId, "\n\n\nThere is no last record of success\n\n\n")
		return
	}

	//通过上一次的成功部署记录中的MirrorList字段找到对应回滚任务所需要的镜像
	imageInfo, imageErr := models.Mirror{}.GetById(lastSuccDeployInfo.MirrorList)
	if imageErr != nil {
		log.Errorf("read mirror sql error: OrderType[%d] , DataId[%d], ErrorReason[%s]\n", RollbackType, lastSuccDeployInfo.MirrorList, imageErr.Error())
		return
	}
	if imageInfo == nil {
		log.Errorf("read mirror sql error: OrderType[%d] , DataId[%d], ErrorReason[no id in sql]\n", RollbackType, lastSuccDeployInfo.MirrorList)
		return
	}

	imageName := fmt.Sprintf("%s:%s", imageInfo.MirrorName, imageInfo.MirrorVersion)

	// 解析机器列表
	var machineList deploy.MachineListJson
	err := json.Unmarshal([]byte(deployInfo.HostList), &machineList)
	if err != nil {
		log.Errorf("deploy hostlist json error: OrderType[%d] , DataId[%d], hostlist[%s], ErrorReason[%s]\n", RollbackType, dataId, deployInfo.HostList, err.Error())
		return
	}

	machineNum, machineSuccNum := deploy.GetMachineSum(machineList)
	fmt.Println(machineNum, machineSuccNum)

	// 解析上一次成功的部署记录中的dockerconfig，用来做构建镜像
	var dockerConf deploy.CreateContainerJson
	err = json.Unmarshal([]byte(lastSuccDeployInfo.DockerConfig), &dockerConf)
	if err != nil {
		log.Errorf("deploy docker config json error: OrderType[%d] , DataId[%d], docker config[%s], ErrorReason[%s]\n", RollbackType, dataId, lastSuccDeployInfo.DockerConfig, err.Error())
		return
	}

	//解析，讲输入的cmd字符串以空格切分为字符串数组格式
	var cmdArr []string
	if len(dockerConf.Cmd) > 0 {
		cmdArr = strings.Fields(dockerConf.Cmd)
	}

	// 部署操作
	for stageId, stage := range machineList.Stage {
		// 找到所有处于非部署等待状态的阶段，全部进行回滚
		if stage.StageStatus != deploy.STAGE_WAIT {
			// 写入当前处理的阶段
			machineList.StageNum = stageId + 1
			hostList, err := json.Marshal(machineList)
			if err != nil {
				log.Errorf("json marshal error: errReason[%s]\n", err.Error())
				return
			}
			deployInfo.HostList = string(hostList)
			deployErr = models.Deploy{}.Update(deployInfo)
			if deployErr != nil {
				log.Errorf("rewrite deploy record sql error: OrderType[%d] , DataId[%d], ErrorReason[%s]\n", RollbackType, dataId, deployErr.Error())
				return
			}

			for machineId, machine := range stage.Machine {
				// 获取机器信息
				machineInfo, err := models.Host{}.GetById(machine.Id)
				if err != nil {
					log.Errorf("get host sql error: machineId[%d], ErrReason[%s]\n", machine.Id, err.Error())
					//machineList.Stage[stageId].Machine[machineId].MachineStatus = MACHINE_ERR'
					//machineList.Stage[stageId].StageStatus = STAGE_ERR
					//RewriteDeployHostList(dataId, stageId, STAGE_ERR, machineId, MACHINE_ERR)
					return
				}
				if machineInfo == nil {
					log.Errorf("get host sql error: machineId[%d], ErrReason[no id in sql]\n", machine.Id)
					return
				}

				if machine.MachineStatus != deploy.MACHINE_WAIT {
					//处理到某一台机器时先向打印其host
					logStrAdd := fmt.Sprintln("\n\n\n", machineInfo.Ip, "\n\n\n")
					deploy.RewriteDeployLog(deployInfo.Id, logStrAdd)

					//获取对应主机上的容器信息
					containerList, err := docker.ListContainers(machineInfo.Ip)
					if err != nil {
						log.Errorf("get listContainers err: Ip[%s], ErrReason[%s]\n", machineInfo.Ip, err.Error())
						deploy.RewriteDeployHostList(deployInfo.Id, stageId, deploy.STAGE_ERR, machineId, deploy.MACHINE_ERR, err.Error(), -1)
						return
					}
					//获取的容器名会加'/'作为前缀，需加上
					name := fmt.Sprintf("/%s", serviceInfo.ServiceName)
					for _, containerInfo := range containerList {
						if name == containerInfo.Names[0] {
							fmt.Println(name)
							err = docker.StopContainer(machineInfo.Ip, serviceInfo.ServiceName)
							if err != nil {
								deploy.RewriteDeployHostList(deployInfo.Id, stageId, deploy.STAGE_ERR, machineId, deploy.MACHINE_ERR, err.Error(), -1)
								return
							}
							err = docker.RemoveContainer(machineInfo.Ip, serviceInfo.ServiceName, true, false, false)
							if err != nil {
								deploy.RewriteDeployHostList(deployInfo.Id, stageId, deploy.STAGE_ERR, machineId, deploy.MACHINE_ERR, err.Error(), -1)
								return
							}
						}
					}
					//docker.StopContainer(machineInfo.Ip, service.ServiceName)
					//docker.RemoveContainer(machineInfo.Ip, service.ServiceName, true, false, false)
					//fmt.Println(containerList)

					//进行到pull容器的步骤
					//deploy.RewriteDeployLog(deployInfo.Id, "\npull images\n")
					//deploy.RewriteDeployStep(deployInfo.Id, stageId, machineId, deploy.STEP_PULL)
					logStrAdd, err = docker.PullImage(machineInfo.Ip, imageName)
					deploy.RewriteDeployLog(deployInfo.Id, logStrAdd)
					if err != nil {
						deploy.RewriteDeployHostList(deployInfo.Id, stageId, deploy.STAGE_ERR, machineId, deploy.MACHINE_ERR, err.Error(), -1)
						return
					}

					//进行到创建容器的步骤
					//deploy.RewriteDeployStep(deployInfo.Id, stageId, machineId, deploy.STEP_CREATE)

					createParam := docker.CreateContainerConf{
						Host:        machineInfo.Ip,
						ServiceName: serviceInfo.ServiceName,
						Image:       imageName,
						HostName:    dockerConf.HostName,
						Volume:      dockerConf.Volume,
						Expose:      dockerConf.Expose,
						HostList:    dockerConf.HostList,
						WorkDir:     dockerConf.WorkerDir,
						Env:         dockerConf.Env,
						Cmd:         cmdArr,
						Dns:         dockerConf.Dns,
					}

					_, err = docker.CreateContainer(createParam)
					if err != nil {
						deploy.RewriteDeployHostList(deployInfo.Id, stageId, deploy.STAGE_ERR, machineId, deploy.MACHINE_ERR, err.Error(), -1)
						return
					}
					//deploy.RewriteDeployLog(deployInfo.Id, fmt.Sprintf("\ncreate container: %s\n", containerId))

					//进行启动容器的步骤
					//deploy.RewriteDeployLog(deployInfo.Id, "\nstart container\n")
					//deploy.RewriteDeployStep(deployInfo.Id, stageId, machineId, deploy.STEP_START)
					err = docker.StartContainer(machineInfo.Ip, serviceInfo.ServiceName)
					if err != nil {
						deploy.RewriteDeployHostList(deployInfo.Id, stageId, deploy.STAGE_ERR, machineId, deploy.MACHINE_ERR, err.Error(), -1)
						return
					}
					//deploy.RewriteDeployLog(deployInfo.Id, "\nstart container succ\n")
				}

				//machineSuccNum++
				//progessStatus := machineSuccNum * 100 / machineNum
				//deploy.RewriteDeployHostList(deployInfo.Id, stageId, deploy.STAGE_UNUSE, machineId, deploy.MACHINE_SUCC, "", progessStatus)
				deploy.RewriteDeployHostList(deployInfo.Id, stageId, deploy.STAGE_UNUSE, machineId, deploy.MACHINE_BACK, "", -1)
			}
			deploy.RewriteDeployHostList(deployInfo.Id, stageId, deploy.STAGE_BACK, -1, -1, "", -1)
			break
		}
	}

}
