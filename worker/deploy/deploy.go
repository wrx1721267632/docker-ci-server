/*
@Time : 18-4-20 下午4:57 
@Author : wangruixin
@File : deploy.go
*/

package deploy

import "github.com/wrxcode/deploy-server/docker"

func Deploy() {
	docker.DockerRun()
}