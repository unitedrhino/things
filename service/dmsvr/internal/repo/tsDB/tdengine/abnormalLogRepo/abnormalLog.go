package abnormalLogRepo

import (
	"fmt"
	"gitee.com/unitedrhino/share/clients"
	"gitee.com/unitedrhino/share/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"sync"
)

type AbnormalLogRepo struct {
	t    *clients.Td
	once sync.Once
}

func NewAbnormalLogRepo(dataSource conf.TSDB) *AbnormalLogRepo {
	td, err := clients.NewTDengine(dataSource)
	if err != nil {
		logx.Error("NewTDengine err", err)
		os.Exit(-1)
	}
	return &AbnormalLogRepo{t: td}
}

func (s *AbnormalLogRepo) GetLogStableName() string {
	return fmt.Sprintf("`model_common_abnormal_log`")
}

func (s *AbnormalLogRepo) GetLogTableName(productID, deviceName string) string {
	return fmt.Sprintf("`device_abnormal_log_%s_%s`", productID, deviceName)
}
