package script

import (
	"encoding/json"
	"testing"

	"flag"
	"fmt"

	"github.com/wrxcode/deploy-server/common"
)

func TestCheckContainer(t *testing.T) {
	machineList := MachineListJson{}
	machineList.Stage = []StageJson{
		{
			Machine: []MachineJson{
				{
					Id:              1,
					ContainerStatus: "",
				},
				{
					Id:              2,
					ContainerStatus: "",
				},
			},
		},
	}

	hostList, err := json.Marshal(machineList)
	if err != nil {
		fmt.Printf("json marshal error: errReason[%s]\n", err.Error())
	}
	fmt.Println(string(hostList))

}

func TestCheckContainer2(t *testing.T) {
	cfgFile := flag.String("c", "../cfg/cfg.toml.debug", "set config file")
	flag.Parse()
	common.Init(*cfgFile)

	CheckContainer()
}
