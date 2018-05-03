package deploy

import (
	"encoding/json"
	"fmt"
	"testing"
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

}
