package sdkLogRepo

import (
	"fmt"
	"gitee.com/unitedrhino/share/clients"
	"gitee.com/unitedrhino/share/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

type SDKLogRepo struct {
	t *clients.Td
	SDKLogStore
}

func NewSDKLogRepo(dataSource conf.TSDB) *SDKLogRepo {
	td, err := clients.NewTDengine(dataSource)
	if err != nil {
		logx.Error("tdengine 初始化错误 err", err)
		os.Exit(-1)
	}
	return &SDKLogRepo{t: td}
}

type SDKLogStore struct {
}

func (s *SDKLogStore) GetSDKLogStableName() string {
	return fmt.Sprintf("`model_common_sdklog`")
}

func (s *SDKLogStore) GetSDKLogTableName(productID, deviceName string) string {
	return fmt.Sprintf("`device_sdklog_%s_%s`", productID, deviceName)
}
