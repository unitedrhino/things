package devicemanagelogic

import (
	"database/sql"
	"encoding/json"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/i-Things/things/shared/domain/schema"
	mysql "github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/spf13/cast"
)

func GetNullTime(time sql.NullTime) int64 {
	if time.Valid == false {
		return 0
	}
	return time.Time.Unix()
}

func ToProductSchema(pt *schema.SchemaInfo) *dm.ProductSchema {
	return &dm.ProductSchema{
		CreatedTime: pt.CreatedTime.Unix(),
		ProductID:   pt.ProductID,
		Schema:      pt.Template,
	}
}

func ToDeviceInfo(di *mysql.DeviceInfo) *dm.DeviceInfo {
	var (
		tags map[string]string
	)

	if di.Tags.String != "" {
		_ = json.Unmarshal([]byte(di.Tags.String), &tags)
	}
	return &dm.DeviceInfo{
		Version:     &wrappers.StringValue{Value: di.Version},
		LogLevel:    di.LogLevel,
		Cert:        di.Cert,
		ProductID:   di.ProductID,
		DeviceName:  di.DeviceName,
		CreatedTime: di.CreatedTime.Unix(),
		FirstLogin:  GetNullTime(di.FirstLogin),
		LastLogin:   GetNullTime(di.LastLogin),
		Secret:      di.Secret,
		IsOnline:    cast.ToBool(di.IsOnline),
		Tags:        tags,
	}
}
