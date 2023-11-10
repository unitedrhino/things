package main

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/disvr/didirect"
)

func main() {
	defer utils.Recover(context.Background())
	svcCtx := didirect.GetSvcCtx()
	didirect.Run(svcCtx)
}
