package deploy

import (
	"encoding/json"
	"flag"
	"fmt"
	"testing"

	"strings"

	"github.com/wrxcode/deploy-server/common"
)

func TestDeploy(t *testing.T) {
	//json_str := "{\"stage\":[{\"stage_status\":1,\"machine\":[{\"id\":111, \"machine_status\":0, \"step\":\"\"},{\"id\":112, \"machine_status\":0, \"step\":\"\"}]},{\"stage_status\":0,\"machine\":[{\"id\":111, \"machine_status\":0, \"step\":\"\"},{\"id\":112, \"machine_status\":0, \"step\":\"\"}],\"stage_status\":0,\"machine\":[{\"id\":111, \"machine_status\":0, \"step\":\"\"},{\"id\":112, \"machine_status\":0, \"step\":\"\"}]}],\"stage_num\":111,\"progress_status\":30}"
	//json_str := "{\"stage\":[],\"stage_num\":111,\"progress_status\":30}"

	//var dat MachineListJson
	//err := json.Unmarshal([]byte(json_str), &dat)
	//fmt.Println("data: ", dat)
	//fmt.Println("size: ", len(dat.Stage))
	//fmt.Println("err: ", err)
	//fmt.Println(json_str)
	//str, err := json.Marshal(dat)
	//fmt.Println(string(str), "    err:", err)
	hostList := MachineListJson{}
	hostList.StageNum = 0
	hostList.ProgressStatus = 0
	hostList.Stage = []StageJson{
		{
			StageStatus: STAGE_DOING,
			Machine: []MachineJson{
				{
					Id:            1,
					Step:          "",
					MachineStatus: MACHINE_DOING,
				},
				{
					Id:            2,
					Step:          "",
					MachineStatus: MACHINE_DOING,
				},
			},
		},
	}

	str, err := json.Marshal(hostList)
	fmt.Println(string(str), "    err:", err)

	deploy_config := CreateContainerJson{
		WorkerDir: "",
		HostName:  "www.test.com",
		HostList:  []string{"www.baidu.com:127.0.0.1"},
		Env:       []string{},
		Volume:    []string{"/root/registry:/var/lib/registry"},
		Dns:       []string{"127.0.0.1", "114.114.114.114"},
		Expose:    []string{"9003:80"},
		//Cmd:       []string{},
	}

	docker_str, err := json.Marshal(deploy_config)
	fmt.Println(string(docker_str), "    err:", err)

}

func TestDeploy2(t *testing.T) {
	cfgFile := flag.String("c", "../../cfg/cfg.toml.debug", "set config file")
	flag.Parse()
	common.Init(*cfgFile)
	Deploy(1)
}

func TestDeploy3(t *testing.T) {
	//var cmdArr []string
	fmt.Println(strings.Fields("a  b   c    d"))
}
