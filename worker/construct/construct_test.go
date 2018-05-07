package construct

import (
	"flag"
	"testing"

	"github.com/wrxcode/deploy-server/common"
)

func TestDockerPull(t *testing.T) {
	cfgFile := flag.String("c", "../../cfg/cfg.toml.debug", "set config file")
	flag.Parse()
	common.Init(*cfgFile)

	ConstructImage(20)
}
