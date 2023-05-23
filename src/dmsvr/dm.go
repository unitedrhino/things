// 设备管理模块-dmsvr
package main

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/dmdirect"
	_ "net/http/pprof"
)

func main() {
	defer utils.Recover(context.Background())
	svcCtx := dmdirect.GetSvcCtx()
	dmdirect.RunServer(svcCtx)
}
