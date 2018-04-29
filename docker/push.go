package docker

import (
	"fmt"
	"github.com/wrxcode/deploy-server/common/g"
	"github.com/pkg/errors"
	"os/exec"
	"bytes"
)

func DockerPush(projectName string, tag string) (string, error) {
	var repo string
	if g.Conf().Repo.IsHost == 1 {
		repo = g.Conf().Repo.Host
	} else if g.Conf().Repo.IsIp == 1 {
		repo = g.Conf().Repo.Ip
	}
	if repo == "" {
		return "", errors.Errorf("config file error:[repo]")
	}

	repo = fmt.Sprintf("%s:%s", repo, g.Conf().Repo.Port)

	//拼接镜像名与私有仓库名，方便docker push使用
	repoData := fmt.Sprintf("%s/%s:%s", repo, projectName, tag)

	cmd := exec.Command("docker", "push", repoData)
	outbuf := new(bytes.Buffer)
	errbuf := new(bytes.Buffer)
	cmd.Stdout = outbuf
	cmd.Stderr = errbuf
	err := cmd.Run()
	if err != nil {
		return errbuf.String(), err
	}
	if errbuf.String() != "" {
		return errbuf.String(), errors.Errorf("docker pull exec error")
	}

	return outbuf.String(), nil
}