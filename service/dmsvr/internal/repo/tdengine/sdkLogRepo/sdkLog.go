package sdkLogRepo

import (
	"gitee.com/unitedrhino/share/clients"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/stores"
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
