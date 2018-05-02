package docker

import (
	"testing"
	"fmt"
)

type CreateContainerJson struct{
	WorkerDir 	string		`json:"workdir"`
	HostName 	string		`json:"hostname"`
	HostList 	[]string 	`json:"hostlist"`
	Env 		[]string	`json:"env"`
	Volume 		[]string	`json:"Volume"`
	Dns 		[]string	`json:"dns"`
	Expose 		[]string	`json:"expose"`
	Cmd 		[]string	`json:"cmd"`
}

func TestCreateContainer(t *testing.T) {
	//str, err := CreateContainer("xupt3.fightcoder.com:9005/nginx:4a5ee6a7dbb494909d12a9d6bee8a791289fc240", cmd)
	param := CreateContainerConf {
		host: 			"tcp://222.24.63.117:9000",
		projectName: 	"registry",
		image: 			"hub.c.163.com/library/registry",
		hostName: 		"www.test.com",
		volume:			[]string{"/root/registry/:/var/lib/registry"},
		expose: 		[]string{"9005:5000"},
		hostList: 		[]string{"www.baidu.com:127.0.0.1"},
	}
	str, err := CreateContainer(param)
	fmt.Println("create container ID:", str, "   err:", err)
}
