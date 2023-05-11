//设备管理模块-dmsvr
package main

import (
	"github.com/i-Things/things/src/dmsvr/dmdirect"
	_ "net/http/pprof"
)

func main() {
	svcCtx := dmdirect.GetSvcCtx()
	dmdirect.RunServer(svcCtx)
}
