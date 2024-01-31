// 设备管理模块-dmsvr
package main

import (
	"context"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/dmsvr/dmdirect"
	_ "net/http/pprof"
)

func main() {
	defer utils.Recover(context.Background())
	svcCtx := dmdirect.GetSvcCtx()
	dmdirect.Run(svcCtx)
}
