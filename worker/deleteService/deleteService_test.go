package delete

import (
	"flag"
	"testing"

	"github.com/wrxcode/deploy-server/common"
)

func TestDeleteService(t *testing.T) {
	cfgFile := flag.String("c", "../../cfg/cfg.toml.debug", "set config file")
	flag.Parse()
	common.Init(*cfgFile)
	DeleteService(2)
}
