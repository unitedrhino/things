// api网关接口代理模块-apisvr
package main

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/apidirect"
	"github.com/zeromicro/go-zero/core/logx"
	_ "github.com/zeromicro/go-zero/core/proc" //开启pprof采集 https://mp.weixin.qq.com/s/yYFM3YyBbOia3qah3eRVQA
)

func main() {
	defer utils.Recover(context.Background())
	logx.DisableStat()
	apiCtx := apidirect.NewApi(apidirect.ApiCtx{})
	apiCtx.Server.PrintRoutes()
	fmt.Printf("Starting apiSvr at %s:%d...\n", apiCtx.Svc.Config.Host, apiCtx.Svc.Config.Port)
	apiCtx.Server.Start()
	defer apiCtx.Server.Stop()
}
