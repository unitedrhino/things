package sdkLogRepo

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/store"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsg"
	"github.com/zeromicro/go-zero/core/logx"
)

func (d SDKLogRepo) GetDeviceSDKLog(ctx context.Context,
	productID, deviceName string, page def.PageInfo2) ([]*deviceMsg.SDKLog, error) {
	sqSql := sq.Select("*").From(getSDKLogStableName(productID)).
		Where("`device_name`=?", deviceName).OrderBy("`ts` desc")
	sqSql = page.FmtSql(sqSql)
	sqlStr, value, err := sqSql.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := d.t.Query(sqlStr, value...)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		} else {
			return []*deviceMsg.SDKLog{}, nil
		}
	}
	var datas []map[string]any
	store.Scan(rows, &datas)
	retLogs := make([]*deviceMsg.SDKLog, 0, len(datas))
	for _, v := range datas {
		retLogs = append(retLogs, ToDeviceSDKLog(productID, v))
	}
	return retLogs, nil
}

func (d SDKLogRepo) Insert(ctx context.Context, data *deviceMsg.SDKLog) error {
	sql := fmt.Sprintf("insert into %s using %s tags('%s')(`ts`, `content`,`log_level`,`client_token`) values (?,?,?,?);",
		getSDKLogTableName(data.ProductID, data.DeviceName), getSDKLogStableName(data.ProductID), data.DeviceName)
	if _, err := d.t.Exec(sql, data.Timestamp, data.Content, data.LogLevel, data.ClientToken); err != nil {
		logx.WithContext(ctx).Errorf(
			sql+"%s|EventTable|productID:%v,deviceName:%v,err:%v",
			utils.FuncName(), data.ProductID, data.DeviceName, err)
		return err
	}
	return nil
}
