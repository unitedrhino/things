package statusLogRepo

import (
	"fmt"
	"gitee.com/unitedrhino/share/clients"
	"gitee.com/unitedrhino/share/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"sync"
)

type StatusLogRepo struct {
	t    *clients.Td
	once sync.Once
}

func NewStatusLogRepo(dataSource conf.TSDB) *StatusLogRepo {
	td, err := clients.NewTDengine(dataSource)
	if err != nil {
		logx.Error("NewTDengine err", err)
		os.Exit(-1)
	}
	return &StatusLogRepo{t: td}
}

func (s *StatusLogRepo) GetLogStableName() string {
	return fmt.Sprintf("`model_common_status_log`")
}

func (s *StatusLogRepo) GetLogTableName(productID, deviceName string) string {
	return fmt.Sprintf("`device_status_log_%s_%s`", productID, deviceName)
}
