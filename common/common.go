package common

import (
	"github.com/wrxcode/deploy-server/common/g"
	"github.com/wrxcode/deploy-server/common/store"
)

func Init(cfgFile string) {
	g.LoadConfig(cfgFile)
	g.InitLog()
	store.InitMysql()
	store.InitRedis()
}

func Close() {
	store.CloseMysql()
	store.CloseRedis()
}
