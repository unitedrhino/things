package main

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/vidsip/sipdirect"
)

func main() {
	defer utils.Recover(context.Background())
	svcCtx := sipdirect.GetSvcCtx()
	sipdirect.Run(svcCtx)
}
