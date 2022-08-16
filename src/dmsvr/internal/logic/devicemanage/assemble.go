package devicemanagelogic

import (
	"encoding/json"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/utils"
	mysql "github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
)

func ToDeviceInfo(di *mysql.DeviceInfo) *dm.DeviceInfo {
	var (
		tags map[string]string
	)

	if di.Tags.String != "" {
		_ = json.Unmarshal([]byte(di.Tags.String), &tags)
	}
	if di.IsOnline == def.Unknown {
		di.IsOnline = def.OffLine
	}
	if di.LogLevel == def.Unknown {
		di.LogLevel = def.LogClose
	}
	return &dm.DeviceInfo{
		Version:     &wrappers.StringValue{Value: di.Version},
		LogLevel:    di.LogLevel,
		Cert:        di.Cert,
		ProductID:   di.ProductID,
		DeviceName:  di.DeviceName,
		CreatedTime: di.CreatedTime.Unix(),
		FirstLogin:  utils.GetNullTime(di.FirstLogin),
		LastLogin:   utils.GetNullTime(di.LastLogin),
		Secret:      di.Secret,
		IsOnline:    int32(di.IsOnline),
		Tags:        tags,
	}
}
