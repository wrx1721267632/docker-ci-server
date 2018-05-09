package delete

import (
	"encoding/json"

	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/wrxcode/deploy-server/docker"
	"github.com/wrxcode/deploy-server/models"
)

type MachineJson struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	ContainerStatus string `json:"container_status"`
}

type StageJson struct {
	Machine []MachineJson `json:"machine"`
}

//服务表中机器列表的json
type MachineListJson struct {
	Stage []StageJson `json:"stage"`
}

const (
	DeleteType = 3
)

func DeleteService(dataId int64) {
	//通过dataID找到对应服务信息
	service, serviceErr := models.Service{}.GetById(dataId)
	if serviceErr != nil {
		log.Errorf("read service sql error: OrderType[%d] , DataId[%d], ErrorReason[%s]\n", DeleteType, dataId, serviceErr.Error())
		return
	}
	if service == nil {
		log.Errorf("read service sql error: OrderType[%d] , DataId[%d], ErrorReason[no id in sql]\n", DeleteType, dataId)
		return
	}

	// 解析机器列表
	var machineList MachineListJson
	err := json.Unmarshal([]byte(service.HostList), &machineList)
	if err != nil {
		log.Errorf("service hostlist json error: OrderType[%d] , DataId[%d], hostlist[%s], ErrorReason[%s]\n", DeleteType, dataId, service.HostList, err.Error())
		return
	}

	for _, stage := range machineList.Stage {
		for _, machine := range stage.Machine {
			// 获取机器信息
			machineInfo, err := models.Host{}.GetById(machine.Id)
			if err != nil {
				log.Errorf("get host sql error: machineId[%d], ErrReason[%s]\n", machine.Id, err.Error())
				//machineList.Stage[stageId].Machine[machineId].MachineStatus = MACHINE_ERR
				//machineList.Stage[stageId].StageStatus = STAGE_ERR
				//RewriteDeployHostList(dataId, stageId, STAGE_ERR, machineId, MACHINE_ERR)
				return
			}
			if machineInfo == nil {
				log.Errorf("get host sql error: machineId[%d], ErrReason[no id in sql]\n", machine.Id)
				return
			}

			//处理到某一台机器时先向打印其host
			//logStrAdd := fmt.Sprintln("\n\n\n", machineInfo.Ip, "\n\n\n")
			//RewriteDeployLog(deploy.Id, logStrAdd)

			//获取对应主机上的容器信息
			containerList, err := docker.ListContainers(machineInfo.Ip)
			if err != nil {
				log.Errorf("get listContainers err: Ip[%s], ErrReason[%s]\n", machineInfo.Ip, err.Error())
				//RewriteDeployHostList(deploy.Id, stageId, STAGE_ERR, machineId, MACHINE_ERR, err.Error(), -1)
				return
			}
			//获取的容器名会加'/'作为前缀，需加上
			name := fmt.Sprintf("/%s", service.ServiceName)
			for _, containerInfo := range containerList {
				if name == containerInfo.Names[0] {
					err = docker.StopContainer(machineInfo.Ip, service.ServiceName)
					if err != nil {
						log.Errorf("stop Containers err: Ip[%s], ErrReason[%s]\n", machineInfo.Ip, err.Error())
						//RewriteDeployHostList(deploy.Id, stageId, STAGE_ERR, machineId, MACHINE_ERR, err.Error(), -1)
						return
					}
					err = docker.RemoveContainer(machineInfo.Ip, service.ServiceName, true, false, false)
					if err != nil {
						log.Errorf("remove Containers err: Ip[%s], ErrReason[%s]\n", machineInfo.Ip, err.Error())
						//RewriteDeployHostList(deploy.Id, stageId, STAGE_ERR, machineId, MACHINE_ERR, err.Error(), -1)
						return
					}
				}
			}
		}
	}

}
