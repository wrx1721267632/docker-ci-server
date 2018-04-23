package models

import (
	"sync"

	"github.com/wrxcode/deploy-server/common"
)

var once sync.Once

func InitAllInTest() {
	once.Do(func() {
		common.Init("../cfg/cfg.toml.debug")
	})
}
