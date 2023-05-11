//系统管理模块-syssvr
package main

import (
	"github.com/i-Things/things/src/syssvr/sysdirect"
)

func main() {
	svcCtx := sysdirect.GetSvcCtx()
	sysdirect.RunServer(svcCtx)
}
