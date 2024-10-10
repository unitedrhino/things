package main

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dgsvr/dgdirect"
)

func main() {
	defer utils.Recover(context.Background())
	svcCtx := dgdirect.GetSvcCtx()
	dgdirect.Run(svcCtx)
}
