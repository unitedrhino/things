// 系统管理模块-syssvr
package main

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/sysdirect"
)

func main() {
	defer utils.Recover(context.Background())
	svcCtx := sysdirect.GetSvcCtx()
	sysdirect.RunServer(svcCtx)
}
