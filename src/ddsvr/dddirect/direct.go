package dddirect

import (
	"github.com/i-Things/things/src/ddsvr/internal/config"
	"github.com/i-Things/things/src/ddsvr/internal/startup"
)

type Config = config.Config

func NewDd(config Config) {
	startup.NewDd(config)
}
