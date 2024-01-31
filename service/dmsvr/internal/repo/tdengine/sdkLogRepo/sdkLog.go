package sdkLogRepo

import (
	"gitee.com/i-Things/share/clients"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/stores"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

type SDKLogRepo struct {
	t *clients.Td
	stores.SDKLogStore
}

func NewSDKLogRepo(dataSource conf.TSDB) *SDKLogRepo {
	td, err := clients.NewTDengine(dataSource)
	if err != nil {
		logx.Error("tdengine 初始化错误 err", err)
		os.Exit(-1)
	}
	return &SDKLogRepo{t: td}
}
