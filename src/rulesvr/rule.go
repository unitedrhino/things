// 协议规则引擎模块-rulesvr
package main

import (
	"context"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/rulesvr/ruledirect"
)

func main() {
	defer utils.Recover(context.Background())
	svcCtx := ruledirect.GetSvcCtx()
	ruledirect.Run(svcCtx)
}
