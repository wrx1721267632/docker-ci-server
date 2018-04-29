package docker

import (
	"testing"
	"flag"
	"github.com/wrxcode/deploy-server/common"
	"fmt"
)

func TestDockerPull(t *testing.T) {
	cfgFile := flag.String("c", "../cfg/cfg.toml.debug", "set config file")
	flag.Parse()
	common.Init(*cfgFile)
	/*var gitpath string
	flag.StringVar(&gitpath, "p", "https://github.com/wrxcode/nginx-docker", "Git repository")
	fmt.Printf("test: %s\n", gitpath)
	ret, err := commit.GetCommit(gitpath)
	fmt.Printf("commit: [%s]\nerror: [%v]\n", ret, err)*/
	/*err1 := DockerBuild("https://github.com/wrxcode/nginx-docker", "test", "111")
	fmt.Printf("error: [%v]\n", err1)*/
	//store.InitMysql()

	str, err := DockerPush("test","111")
	fmt.Println(str, "\n-------------------\n", err)
}