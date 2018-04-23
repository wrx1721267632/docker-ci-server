package main

import (
	"flag"
	"github.com/wrxcode/deploy-server/common"
	"github.com/wrxcode/deploy-server/router"
	"github.com/TV4/graceful"
	"net/http"
	"github.com/wrxcode/deploy-server/common/g"
	"github.com/wrxcode/deploy-server/dispatch"
	"fmt"
)

func main() {
	cfgFile := flag.String("c", "cfg/cfg.toml.debug", "set config file")
	flag.Parse()

	common.Init(*cfgFile)
	defer common.Close()

	router := router.GetRouter()

	go dispatch.StartConsume()

	graceful.LogListenAndServe(&http.Server{
		Addr:    fmt.Sprintf(":%d", g.Conf().Run.HTTPPort),
		Handler: router,
	})
}
