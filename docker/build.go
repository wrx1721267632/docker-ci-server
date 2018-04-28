/*
@Time : 18-4-20 下午8:33 
@Author : wangruixin
@File : build.go
*/

package docker

import (
	"os/exec"
	"bytes"
	"strings"
	"github.com/pkg/errors"
)

func DockerBuild(gitPath string, repo string, tag string) (string, error) {
	gitPath = strings.Replace(gitPath, "http://", "", -1)
	gitPath = strings.Replace(gitPath, "https://", "", -1)

	repoData := repo
	repoData += ":"
	repoData += tag

	cmd := exec.Command("docker", "build", "-t", repoData, gitPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	var err_out bytes.Buffer
	cmd.Stderr = &err_out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	if err_out.String() != "" {
		return err_out.String(), errors.Errorf("docker build exec error")
	}

	return out.String(), nil
}

func DockerTag() {

}
