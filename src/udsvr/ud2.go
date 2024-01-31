package main

import (
	"context"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/udsvr/uddirect"
)

func main() {
	defer utils.Recover(context.Background())
	svcCtx := uddirect.GetSvcCtx()
	uddirect.Run(svcCtx)
}
