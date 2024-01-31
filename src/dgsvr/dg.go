package main

import (
	"context"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/dgsvr/dgdirect"
)

func main() {
	defer utils.Recover(context.Background())
	svcCtx := dgdirect.GetSvcCtx()
	dgdirect.Run(svcCtx)
}
