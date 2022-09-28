package main

import (
	"github.com/i-Things/things/src/apisvr/apidirect"
	"github.com/zeromicro/go-zero/core/logx"
)

func main() {
	logx.DisableStat()
	server := apidirect.NewApi(apidirect.ApiCtx{})
	server.Start()
}
