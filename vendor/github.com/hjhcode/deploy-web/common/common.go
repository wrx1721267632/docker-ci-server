package common

import (
	"github.com/hjhcode/deploy-web/common/g"
	"github.com/hjhcode/deploy-web/common/store"
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
