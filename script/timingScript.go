package script

import (
	"encoding/json"

	"fmt"

	"time"

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

// 定时任务获取容器状态
func CheckContainer() {
	for {
		//定时一分钟执行一次
		time.Sleep(10 * time.Second)

		serviceAll, err := models.Service{}.QueryAllService()
		if err != nil {
			log.Error("mysql error: read service QueryAllService")
			panic(err)
		}

		for _, service := range serviceAll {
			//fmt.Println("start")
			// 解析机器列表
			var machineList MachineListJson
			err = json.Unmarshal([]byte(service.HostList), &machineList)
			if err != nil {
				log.Errorf("timing script service hostlist json error: hostlist[%s], ErrorReason[%s]\n", service.HostList, err.Error())
				continue
			}

			for stageId, stage := range machineList.Stage {
				for machineId, machine := range stage.Machine {
					// 获取机器信息
					machineInfo, err := models.Host{}.GetById(machine.Id)
					if err != nil {
						log.Errorf("get host sql error: machineId[%d], ErrReason[%s]\n", machine.Id, err.Error())
						continue
					}
					if machineInfo == nil {
						log.Errorf("get host sql error: machineId[%d], ErrReason[no id in sql]\n", machine.Id)
						continue
					}

					//获取对应主机上的容器信息
					containerList, err := docker.ListContainers(machineInfo.Ip)
					if err != nil {
						log.Errorf("get listContainers err: Ip[%s], ErrReason[%s]\n", machineInfo.Ip, err.Error())
						machineList.Stage[stageId].Machine[machineId].ContainerStatus = err.Error()
						continue
					}

					//获取的容器名会加'/'作为前缀，需加上
					name := fmt.Sprintf("/%s", service.ServiceName)

					flag := false
					for _, containerInfo := range containerList {
						if name == containerInfo.Names[0] {
							machineList.Stage[stageId].Machine[machineId].ContainerStatus = containerInfo.State
							if containerInfo.State != "running" {
								err = docker.StartContainer(machineInfo.Ip, service.ServiceName)
								if err != nil {
									log.Errorf("restart container error: host[%s], service[%s]\n", machineInfo.Ip, service.ServiceName)
									continue
								}
								log.Errorf("restart the container: host[%s], service[%s]\n", machineInfo.Ip, service.ServiceName)
							}
							flag = true
							break
						}
					}

					if flag == false {
						machineList.Stage[stageId].Machine[machineId].ContainerStatus = "There is no container in the machine!"
						log.Errorf("There is no container in the machine: host[%s], service[%s]\n", machineInfo.Ip, service.ServiceName)
					}
				}
			}

			hostList, err := json.Marshal(machineList)
			if err != nil {
				log.Errorf("json marshal error: errReason[%s]\n", err.Error())
				continue
			}

			service.HostList = string(hostList)
			err = models.Service{}.Update(service)
			if err != nil {
				log.Errorf("rewrite service record sql error: ErrorReason[%s]", err)
				continue
			}
			//fmt.Println("end")
		}
	}
}
