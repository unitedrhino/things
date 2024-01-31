package main

import (
	"context"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/vidsip/sipdirect"
)

func main() {
	defer utils.Recover(context.Background())
	svcCtx := sipdirect.GetSvcCtx()
	sipdirect.Run(svcCtx)
}
