package main

import (
	"context"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/service/vidsip/sipdirect"
)

func main() {
	defer utils.Recover(context.Background())
	svcCtx := sipdirect.GetSvcCtx()
	sipdirect.Run(svcCtx)
}
