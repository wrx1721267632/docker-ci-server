package construct

import (
"testing"
"flag"
"github.com/wrxcode/deploy-server/common"
)

func TestDockerPull(t *testing.T) {
	cfgFile := flag.String("c", "../../cfg/cfg.toml.debug", "set config file")
	flag.Parse()
	common.Init(*cfgFile)

	ContructImage(2)
}