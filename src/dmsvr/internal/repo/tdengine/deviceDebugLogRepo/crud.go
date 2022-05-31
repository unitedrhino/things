package deviceDebugLogRepo

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/store"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/domain/device"
	"github.com/zeromicro/go-zero/core/logx"
)

func (d DeviceDebugLogRepo) GetDeviceDebugLog(ctx context.Context, productID, deviceName string, page def.PageInfo2) ([]*device.DebugLog, error) {
	sql := sq.Select("*").From(getDebugLogStableName(productID)).
		Where("`device_name`=?", deviceName).OrderBy("`ts` desc")
	sql = page.FmtSql(sql)
	sqlStr, value, err := sql.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := d.t.Query(sqlStr, value...)
	if err != nil {
		return nil, err
	}
	var datas []map[string]interface{}
	store.Scan(rows, &datas)
	retLogs := make([]*device.DebugLog, 0, len(datas))
	for _, v := range datas {
		retLogs = append(retLogs, ToDeviceDebugLog(productID, v))
	}
	return retLogs, nil
}

func (d DeviceDebugLogRepo) Insert(ctx context.Context, data *device.DebugLog) error {
	sql := fmt.Sprintf("insert into %s using %s tags('%s')(`ts`, `content`,`log_level`,"+
		" `request_id`, `trance_id`, `result_type`) values (?,?,?,?,?,?);",
		getDebugLogTableName(data.ProductID, data.DeviceName), getDebugLogStableName(data.ProductID), data.DeviceName)
	if _, err := d.t.Exec(sql, data.Timestamp, data.Content, data.LogLevel, data.RequestID, data.TranceID, data.ResultType); err != nil {
		logx.WithContext(ctx).Errorf(
			sql+"%s|EventTable|productID:%v,deviceName:%v,err:%v",
			utils.FuncName(), data.ProductID, data.DeviceName, err)
		return err
	}
	return nil
}
