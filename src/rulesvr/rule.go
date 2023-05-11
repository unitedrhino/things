package main

import (
	"github.com/i-Things/things/src/rulesvr/ruledirect"
)

func main() {
	svcCtx := ruledirect.GetSvcCtx()
	ruledirect.RunServer(svcCtx)
}
