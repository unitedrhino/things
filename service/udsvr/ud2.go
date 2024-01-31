package main

import (
	"context"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/udsvr/uddirect"
)

func main() {
	defer utils.Recover(context.Background())
	svcCtx := uddirect.GetSvcCtx()
	uddirect.Run(svcCtx)
}
