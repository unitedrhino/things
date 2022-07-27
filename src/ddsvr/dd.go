package main

import (
	"flag"
	"github.com/i-Things/things/src/ddsvr/internal/config"
	"github.com/i-Things/things/src/ddsvr/internal/startup"
	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/dd.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	startup.NewDd(c)
}
