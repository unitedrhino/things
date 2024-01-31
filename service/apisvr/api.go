// api网关接口代理模块-apisvr
package main

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/service/apisvr/apidirect"
	"github.com/zeromicro/go-zero/core/logx"
	_ "github.com/zeromicro/go-zero/core/proc" //开启pprof采集 https://mp.weixin.qq.com/s/yYFM3YyBbOia3qah3eRVQA
)

func main() {
	defer utils.Recover(context.Background())
	logx.DisableStat()
	apiCtx := apidirect.NewApi(apidirect.ApiCtx{})
	apiCtx.Server.PrintRoutes()
	fmt.Printf("Starting apiSvr at %s:%d...\n", apiCtx.SvcCtx.Config.Host, apiCtx.SvcCtx.Config.Port)
	defer apiCtx.Server.Stop()
	apiCtx.Server.Start()
}
