package dddirect

import (
	"github.com/i-Things/things/src/ddsvr/internal/config"
	"github.com/i-Things/things/src/ddsvr/internal/startup"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

type Config = config.Config

var (
	c config.Config
)

func NewDd() {
	conf.MustLoad("etc/dd.yaml", &c)
	startup.NewDd(c)
	logx.Infof("enabled ddsvr")
}
