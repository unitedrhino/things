// 协议规则引擎模块-rulesvr
package main

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/rulesvr/ruledirect"
)

func main() {
	defer utils.Recover(context.Background())
	svcCtx := ruledirect.GetSvcCtx()
	ruledirect.RunServer(svcCtx)
}
