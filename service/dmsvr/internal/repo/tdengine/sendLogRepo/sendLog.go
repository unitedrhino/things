package sendLogRepo

import (
	"fmt"
	"gitee.com/i-Things/share/clients"
	"gitee.com/i-Things/share/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"sync"
)

type SendLogRepo struct {
	t    *clients.Td
	once sync.Once
}

func NewSendLogRepo(dataSource conf.TSDB) *SendLogRepo {
	td, err := clients.NewTDengine(dataSource)
	if err != nil {
		logx.Error("NewTDengine err", err)
		os.Exit(-1)
	}
	return &SendLogRepo{t: td}
}

func (s *SendLogRepo) GetLogStableName() string {
	return fmt.Sprintf("`model_common_send_log`")
}

func (s *SendLogRepo) GetLogTableName(productID, deviceName string) string {
	return fmt.Sprintf("`device_send_log_%s_%s`", productID, deviceName)
}
