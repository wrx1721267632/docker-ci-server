package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/TV4/graceful"
	"github.com/wrxcode/deploy-server/common"
	"github.com/wrxcode/deploy-server/common/g"
	"github.com/wrxcode/deploy-server/dispatch"
	"github.com/wrxcode/deploy-server/router"
)

func main() {
	cfgFile := flag.String("c", "cfg/cfg.toml.debug", "set config file")
	flag.Parse()

	common.Init(*cfgFile)
	defer common.Close()

	router := router.GetRouter()

	go dispatch.StartConsume()

	//go script.CheckContainer()

	graceful.LogListenAndServe(&http.Server{
		Addr:    fmt.Sprintf(":%d", g.Conf().Run.HTTPPort),
		Handler: router,
	})
}
