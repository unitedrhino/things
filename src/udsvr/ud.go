package main

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/udsvr/uddirect"
)

func main() {
	defer utils.Recover(context.Background())
	svcCtx := uddirect.GetSvcCtx()
	uddirect.Run(svcCtx)
}
