package stockdirect

import (
	"github.com/i-Things/things/src/stocksvr/internal/server"
	client "github.com/i-Things/things/src/stocksvr/stockclient"
)

func NewStock(runSvr bool) client.Stock {
	svcCtx := GetSvcCtx()
	if runSvr {
		RunServer(svcCtx)
	}
	return client.NewDirectStock(svcCtx, server.NewStockServer(svcCtx))
}
