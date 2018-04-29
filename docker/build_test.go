package docker

import (
	"testing"
	"flag"
	"github.com/wrxcode/deploy-server/common"
	"fmt"
)

func TestDockerBuild(t *testing.T) {
	cfgFile := flag.String("c", "../cfg/cfg.toml.debug", "set config file")
	flag.Parse()
	common.Init(*cfgFile)
	/*var gitpath string
	flag.StringVar(&gitpath, "p", "https://github.com/wrxcode/nginx-docker", "Git repository")
	fmt.Printf("test: %s\n", gitpath)
	ret, err := commit.GetCommit(gitpath)
	fmt.Printf("commit: [%s]\nerror: [%v]\n", ret, err)*/
	err1 := DockerBuild("https://github.com/wrxcode/deploy-server", "test", "111", 1)
	fmt.Printf("error: [%v]\n", err1)
	//store.InitMysql()


	//con := models.ConstructRecord{1,1,1,1,time.Now().Unix(),time.Now().Unix(), 0,""}
	//id, err := models.ConstructRecord{}.Add(&con)
	//fmt.Println(id, "        ", err)
}
