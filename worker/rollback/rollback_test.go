package rollback

import (
	"flag"
	"testing"

	"github.com/wrxcode/deploy-server/common"
)

func TestRollback(t *testing.T) {
	cfgFile := flag.String("c", "../../cfg/cfg.toml.debug", "set config file")
	flag.Parse()
	common.Init(*cfgFile)

	Rollback(1)
}
