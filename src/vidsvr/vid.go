package main

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsvr/viddirect"
)

func main() {
	defer utils.Recover(context.Background())
	svcCtx := viddirect.GetSvcCtx()
	viddirect.ApiDirectRun()
	viddirect.Run(svcCtx)
}
