/*
@Time : 18-4-20 下午4:57
@Author : wangruixin
@File : deploy.go
*/

package deploy

import (
	"encoding/json"

	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/wrxcode/deploy-server/docker"
	"github.com/wrxcode/deploy-server/models"
)

const (
	ContructType = 1
)

type MachineJson struct {
	Id            int64  `json:"id"`
	Step          string `json:"Step"`
	MachineStatus int    `json:"machine_status"`
}

type StageJson struct {
	StageStatus int           `json:"stage_status"`
	Machine     []MachineJson `json:"machine"`
}

//机器列表的json
type MachineListJson struct {
	Stage          []StageJson `json:"stage"`
	StageNum       int         `json:"stage_num"`
	ProgressStatus int         `json:"progress_status"`
}

//部署表中的dockerconfig 的json格式解析出来的object
type CreateContainerJson struct {
	WorkerDir string   `json:"workdir"`
	HostName  string   `json:"hostname"`
	HostList  []string `json:"hostlist"`
	Env       []string `json:"env"`
	Volume    []string `json:"Volume"`
	Dns       []string `json:"dns"`
	Expose    []string `json:"expose"`
	Cmd       []string `json:"cmd"`
}

const (
	STEP_PULL   = "pulling image"
	STEP_CREATE = "create container"
	STEP_START  = "start container"
)

const (
	STAGE_WAIT  = 0
	STAGE_DOING = 1
	STAGE_SUCC  = 2
	STAGE_ERR   = 3
	STAGE_UNUSE = -1 //不修改stage状态
)

const (
	MACHINE_WAIT  = 0
	MACHINE_DOING = 1
	MACHINE_ERR   = 2
	MACHINE_SKIP  = 3
	MACHINE_SUCC  = 4
)

func Deploy(dataId int64) {
	//通过dataID获取到对应部署日志表的信息
	deploy, deployErr := models.Deploy{}.GetById(dataId)
	if deployErr != nil {
		log.Errorf("read deploy record sql error: OrderType[%d] , DataId[%d], ErrorReason[%s]\n", ContructType, dataId, deployErr.Error())
		return
	}
	if deploy == nil {
		log.Errorf("read deploy record sql error: OrderType[%d] , DataId[%d], ErrorReason[no id in sql]\n", ContructType, dataId)
		return
	}

	service, serviceErr := models.Service{}.GetById(deploy.ServiceId)
	if serviceErr != nil {
		log.Errorf("read service sql error: OrderType[%d] , DataId[%d], ErrorReason[%s]\n", ContructType, dataId, deployErr.Error())
		return
	}
	if service == nil {
		log.Errorf("read service sql error: OrderType[%d] , DataId[%d], ErrorReason[no id in sql]\n", ContructType, dataId)
		return
	}

	//image, imageErr := models.Mirror{}.GetById(deploy.MirrorId)

	// 解析机器列表
	var machineList MachineListJson
	err := json.Unmarshal([]byte(deploy.HostList), &machineList)
	if err != nil {
		log.Errorf("deploy hostlist json error: OrderType[%d] , DataId[%d], hostlist[%s], ErrorReason[%s]\n", ContructType, dataId, deploy.HostList, err.Error())
		return
	}

	machineNum, machineSuccNum := getMachineSum(machineList)
	fmt.Println(machineNum, machineSuccNum)

	// 解析dockerconfig，用来做构建镜像
	var dockerConf CreateContainerJson
	err = json.Unmarshal([]byte(deploy.DockerConfig), &dockerConf)
	if err != nil {
		log.Errorf("deploy docker config json error: OrderType[%d] , DataId[%d], docker config[%s], ErrorReason[%s]\n", ContructType, dataId, deploy.DockerConfig, err.Error())
		return
	}

	// 部署操作
	for stageId, stage := range machineList.Stage {
		// 找到处于部署中状态的阶段（一个时间只有一个阶段会处于部署中，所以直接最后直接break）
		if stage.StageStatus == STAGE_DOING {
			// 写入当前处理的阶段
			machineList.StageNum = stageId + 1
			hostList, err := json.Marshal(machineList)
			if err != nil {
				log.Errorf("json marshal error: errReason[%s]\n", err.Error())
				return
			}
			deploy.HostList = string(hostList)
			deployErr = models.Deploy{}.Update(deploy)
			if deployErr != nil {
				log.Errorf("rewrite deploy record sql error: OrderType[%d] , DataId[%d], ErrorReason[%s]\n", ContructType, dataId, deployErr.Error())
				return
			}

			if deployErr != nil {
				log.Errorf("read deploy record sql error: DataId[%d], ErrorReason[%s]", deploy.Id, deployErr)
			}

			for machineId, machine := range stage.Machine {
				if machine.MachineStatus == MACHINE_DOING {
					// 获取机器信息
					machineInfo, err := models.Host{}.GetById(machine.Id)
					if err != nil {
						log.Errorf("get host sql error: machineId[%d], ErrReason[%s]\n", machine.Id, err.Error())
						//machineList.Stage[stageId].Machine[machineId].MachineStatus = MACHINE_ERR
						//machineList.Stage[stageId].StageStatus = STAGE_ERR
						//rewriteDeployHostList(dataId, stageId, STAGE_ERR, machineId, MACHINE_ERR)
						return
					}
					if machineInfo == nil {
						log.Errorf("get host sql error: machineId[%d], ErrReason[no id in sql]\n", machine.Id)
						return
					}

					//处理到某一台机器时先向打印其host
					logStrAdd := fmt.Sprintln("\n\n\n", machineInfo.Ip, "\n\n\n")
					rewriteDeployLog(deploy.Id, logStrAdd)

					//获取对应主机上的容器信息
					containerList, err := docker.ListContainers(machineInfo.Ip)
					if err != nil {
						log.Errorf("get listContainers err: Ip[%s], ErrReason[%s]\n", machineInfo.Ip, err.Error())
						rewriteDeployHostList(deploy.Id, stageId, STAGE_ERR, machineId, MACHINE_ERR, err.Error(), -1)
						return
					}
					//获取的容器名会加'/'作为前缀，需加上
					name := fmt.Sprintf("/", service.ServiceName)
					for _, containerInfo := range containerList {
						if name == containerInfo.Names[0] {
							err = docker.StopContainer(machineInfo.Ip, service.ServiceName)
							if err != nil {
								rewriteDeployHostList(deploy.Id, stageId, STAGE_ERR, machineId, MACHINE_ERR, err.Error(), -1)
								return
							}
							err = docker.RemoveContainer(machineInfo.Ip, service.ServiceName, true, false, false)
							if err != nil {
								rewriteDeployHostList(deploy.Id, stageId, STAGE_ERR, machineId, MACHINE_ERR, err.Error(), -1)
								return
							}
						}
					}
					//docker.StopContainer(machineInfo.Ip, service.ServiceName)
					//docker.RemoveContainer(machineInfo.Ip, service.ServiceName, true, false, false)
					//fmt.Println(containerList)

					//进行到pull容器的步骤
					rewriteDeployStep(deploy.Id, stageId, machineId, STEP_PULL)
					logStrAdd, err = docker.PullImage(machineInfo.Ip, deploy.MirrorList)
					rewriteDeployLog(deploy.Id, logStrAdd)
					if err != nil {
						rewriteDeployHostList(deploy.Id, stageId, STAGE_ERR, machineId, MACHINE_ERR, err.Error(), -1)
						return
					}

					//进行到创建容器的步骤
					rewriteDeployStep(deploy.Id, stageId, machineId, STEP_CREATE)

					createParam := docker.CreateContainerConf{
						Host:        machineInfo.Ip,
						ServiceName: service.ServiceName,
						Image:       deploy.MirrorList,
						HostName:    dockerConf.HostName,
						Volume:      dockerConf.Volume,
						Expose:      dockerConf.Expose,
						HostList:    dockerConf.HostList,
						WorkDir:     dockerConf.WorkerDir,
						Env:         dockerConf.Env,
						Cmd:         dockerConf.Cmd,
						Dns:         dockerConf.Dns,
					}

					_, err = docker.CreateContainer(createParam)
					if err != nil {
						rewriteDeployHostList(deploy.Id, stageId, STAGE_ERR, machineId, MACHINE_ERR, err.Error(), -1)
						return
					}

					rewriteDeployStep(deploy.Id, stageId, machineId, STEP_START)
					err = docker.StartContainer(machineInfo.Ip, service.ServiceName)
					if err != nil {
						rewriteDeployHostList(deploy.Id, stageId, STAGE_ERR, machineId, MACHINE_ERR, err.Error(), -1)
						return
					}
				}

				machineSuccNum++
				progessStatus := machineSuccNum * 100 / machineNum
				rewriteDeployHostList(deploy.Id, stageId, STAGE_UNUSE, machineId, MACHINE_ERR, "", progessStatus)
			}
			rewriteDeployHostList(deploy.Id, stageId, STAGE_SUCC, -1, -1, "", -1)
			break
		}
	}

}

func rewriteDeployLog(deployId int64, logStrAdd string) {
	deploy, deployErr := models.Deploy{}.GetById(deployId)
	if deployErr != nil {
		log.Errorf("read deploy record sql error: OrderType[%d] , DataId[%d], ErrorReason[%s]\n", ContructType, deployId, deployErr.Error())
		return
	}
	if deploy == nil {
		log.Errorf("read deploy record sql error: OrderType[%d] , DataId[%d], ErrorReason[no id in sql]\n", ContructType, deployId)
		return
	}
	deploy.DeployLog += logStrAdd
	deployErr = models.Deploy{}.Update(deploy)
	if deployErr != nil {
		log.Errorf("rewrite deploy record sql error: DataId[%d], ErrorReason[%s]", deploy.Id, deployErr)
	}

}

//修改并回写
func rewriteDeployStep(deployId int64, stageId int, machineId int, step string) {
	deploy, deployErr := models.Deploy{}.GetById(deployId)
	if deployErr != nil {
		log.Errorf("read deploy record sql error: OrderType[%d] , DataId[%d], ErrorReason[%s]\n", ContructType, deployId, deployErr.Error())
		return
	}
	if deploy == nil {
		log.Errorf("read deploy record sql error: OrderType[%d] , DataId[%d], ErrorReason[no id in sql]\n", ContructType, deployId)
		return
	}
	// 解析机器列表
	var machineList MachineListJson
	err := json.Unmarshal([]byte(deploy.HostList), &machineList)
	if err != nil {
		log.Errorf("deploy hostlist json error: OrderType[%d] , DataId[%d], hostlist[%s], ErrorReason[%s]\n", ContructType, deployId, deploy.HostList, err.Error())
		return
	}
	machineList.Stage[stageId].Machine[machineId].Step = step
	hostList, err := json.Marshal(machineList)
	if err != nil {
		log.Errorf("json marshal error: errReason[%s]\n", err.Error())
		return
	}
	deploy.HostList = string(hostList)
	deployErr = models.Deploy{}.Update(deploy)
	if deployErr != nil {
		log.Errorf("rwrite deploy record sql error: DataId[%d], ErrorReason[%s]", deploy.Id, deployErr)
	}
}

// 修改并回写hostlist字段中的stage状态和machine状态
func rewriteDeployHostList(deployId int64, stageId int, stageStatus int, machineId int, machineStatus int, logStrAdd string, progessStatus int) {
	deploy, deployErr := models.Deploy{}.GetById(deployId)
	if deployErr != nil {
		log.Errorf("read deploy record sql error: OrderType[%d] , DataId[%d], ErrorReason[%s]\n", ContructType, deployId, deployErr.Error())
		return
	}
	if deploy == nil {
		log.Errorf("read deploy record sql error: OrderType[%d] , DataId[%d], ErrorReason[no id in sql]\n", ContructType, deployId)
		return
	}
	// 解析机器列表
	var machineList MachineListJson
	err := json.Unmarshal([]byte(deploy.HostList), &machineList)
	if err != nil {
		log.Errorf("deploy hostlist json error: OrderType[%d] , DataId[%d], hostlist[%s], ErrorReason[%s]\n", ContructType, deployId, deploy.HostList, err.Error())
		return
	}

	//修改机器状态
	if machineId != -1 {
		machineList.Stage[stageId].Machine[machineId].MachineStatus = machineStatus
	}

	//修改阶段状态
	if stageStatus != STAGE_UNUSE {
		machineList.Stage[stageId].StageStatus = stageStatus
	}

	if progessStatus != -1 {
		machineList.ProgressStatus = progessStatus
	}

	//讲object转换为json格式字符串
	hostList, err := json.Marshal(machineList)
	if err != nil {
		log.Errorf("json marshal error: errReason[%s]\n", err.Error())
		return
	}
	deploy.HostList = string(hostList)
	deploy.DeployLog += logStrAdd
	deployErr = models.Deploy{}.Update(deploy)
	if deployErr != nil {
		log.Errorf("rewrite deploy record sql error: DataId[%d], ErrorReason[%s]", deploy.Id, deployErr)
	}
}

//获取主机列表中总共的机器个数（用作百分比的分母）和已经部署成功的机器个数（用作百分比分子的基数）
func getMachineSum(machineList MachineListJson) (int, int) {
	//机器总数
	machineNum := 0
	//已经部署成功的机器数
	machineSuccNum := 0
	for _, stage := range machineList.Stage {
		machineNum += len(stage.Machine)
		for _, machine := range stage.Machine {
			if machine.MachineStatus == MACHINE_SUCC {
				machineSuccNum++
			}
		}
	}
	return machineNum, machineSuccNum
}
