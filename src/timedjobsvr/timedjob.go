package main

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/timedjobsvr/timedjobdirect"
)

func main() {
	defer utils.Recover(context.Background())
	ctx := timedjobdirect.GetSvcCtx()
	timedjobdirect.Run(ctx)
}
