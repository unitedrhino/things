package main

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dgsvr/dgdirect"
)

func main() {
	defer utils.Recover(context.Background())
	svcCtx := dgdirect.GetSvcCtx()
	dgdirect.Run(svcCtx)
}
