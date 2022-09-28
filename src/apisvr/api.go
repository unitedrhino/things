package main

import (
	"github.com/i-Things/things/src/apisvr/apidirect"
	"github.com/zeromicro/go-zero/core/logx"
)

func main() {
	logx.DisableStat()
	apiCtx := apidirect.NewApi(apidirect.ApiCtx{})
	apiCtx.Server.Start()
}
