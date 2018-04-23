/*
@Time : 18-4-20 下午6:25 
@Author : wangruixin
@File : rollback.go
*/

package rollback

import "github.com/wrxcode/deploy-server/worker/deploy"

func RollBack() {
	//
	deploy.Deploy()
}
