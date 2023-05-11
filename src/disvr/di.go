package main

import (
	"github.com/i-Things/things/src/disvr/didirect"
)

func main() {
	svcCtx := didirect.GetSvcCtx()
	didirect.RunServer(svcCtx)
}
