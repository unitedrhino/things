package abnormalLogRepo

import (
	"fmt"
	"gitee.com/unitedrhino/share/clients"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceGroup"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"sync"
)

var defaultTags = "`product_id` ,`device_name`,`tenant_code`,`project_id`,`area_id`,`area_id_path` "
var defaultTagDef = "`product_id` BINARY(50),`device_name`  BINARY(50), " +
	"`tenant_code`  BINARY(50),`project_id` BIGINT,`area_id` BIGINT,`area_id_path`  BINARY(50) "

type AbnormalLogRepo struct {
	t            *clients.Td
	once         sync.Once
	groupConfigs []*deviceGroup.GroupDetail
}

func NewAbnormalLogRepo(dataSource conf.TSDB, g []*deviceGroup.GroupDetail) *AbnormalLogRepo {
	td, err := clients.NewTDengine(dataSource)
	if err != nil {
		logx.Error("NewTDengine err", err)
		os.Exit(-1)
	}
	return &AbnormalLogRepo{t: td, groupConfigs: g}
}

func (s *AbnormalLogRepo) GetLogStableName() string {
	return fmt.Sprintf("`model_common_abnormal_log`")
}

func (s *AbnormalLogRepo) GetLogTableName(productID, deviceName string) string {
	return fmt.Sprintf("`device_abnormal_log_%s_%s`", productID, deviceName)
}
