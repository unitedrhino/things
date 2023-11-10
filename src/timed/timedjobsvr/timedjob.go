package main

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/timed/timedjobsvr/timedjobdirect"
)

func main() {
	defer utils.Recover(context.Background())
	ctx := timedjobdirect.GetSvcCtx()
	timedjobdirect.Run(ctx)
}
