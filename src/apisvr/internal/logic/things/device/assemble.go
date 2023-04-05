package device

import (
	"errors"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/spf13/cast"
	"strings"
)

type (
	multiImportCsvRow struct {
		ProductName string //【必填】产品名称
		DeviceName  string //【必填】设备名称 读写
		LogLevel    string //【可选】日志级别:关闭 错误 告警 信息 5调试
		Tags        string //【可选】设备tags
		Position    string //【可选】设备定位,默认百度坐标系
		Address     string //【可选】所在详细地址
	}
)

func DeviceInfoToApi(v *dm.DeviceInfo) *types.DeviceInfo {
	return &types.DeviceInfo{
		ProductID:   v.ProductID,                   //产品id 只读
		DeviceName:  v.DeviceName,                  //设备名称 读写
		CreatedTime: v.CreatedTime,                 //创建时间 只读
		Secret:      v.Secret,                      //设备秘钥 只读
		FirstLogin:  v.FirstLogin,                  //激活时间 只读
		LastLogin:   v.LastLogin,                   //最后上线时间 只读
		Version:     utils.ToNullString(v.Version), // 固件版本  读写
		LogLevel:    v.LogLevel,                    // 日志级别:1)关闭 2)错误 3)告警 4)信息 5)调试  读写
		Cert:        v.Cert,                        // 设备证书  只读
		Tags:        logic.ToTagsType(v.Tags),
		IsOnline:    v.IsOnline, // 在线状态  1离线 2在线 只读
		Address:     &v.Address.Value,
		Position:    &types.Point{Longitude: v.Position.Longitude, Latitude: v.Position.Latitude},
	}
}

func DeviceMultiImportCellToDeviceInfo(cell []string) (info *dm.DeviceInfo, err error) {
	d := &dm.DeviceInfo{}
	m := &multiImportCsvRow{
		ProductName: strings.TrimSpace(utils.SliceIndex(cell, 0, "")),
		DeviceName:  strings.TrimSpace(utils.SliceIndex(cell, 1, "")),
		LogLevel:    strings.TrimSpace(utils.SliceIndex(cell, 2, "")),
		Tags:        strings.TrimSpace(utils.SliceIndex(cell, 3, "")),
		Position:    strings.TrimSpace(utils.SliceIndex(cell, 4, "")),
		Address:     strings.TrimSpace(utils.SliceIndex(cell, 5, "")),
	}

	if m.ProductName == "" {
		return nil, errors.New("缺少必填的产品名称")
	} else {
		d.ProductName = m.ProductName
	}

	if m.DeviceName == "" {
		return nil, errors.New("缺少必填的设备名称")
	} else {
		d.DeviceName = m.DeviceName
	}

	if m.LogLevel != "" {
		if level, ok := def.LogLevelTextToIntMap[m.LogLevel]; !ok {
			return nil, errors.New("设备日志级别格式错误")
		} else {
			d.LogLevel = level
		}
	}

	if m.Tags != "" {
		arr := utils.SplitCutset(m.Tags, ";；\n")
		tagArr := make([]*types.Tag, 0, len(arr))
		for _, item := range arr {
			tagSli := utils.SplitCutset(item, ":：")
			if len(tagSli) == 2 {
				tagArr = append(tagArr, &types.Tag{tagSli[0], tagSli[1]})
			} else {
				return nil, errors.New("设备标签格式错误")
			}
		}
		d.Tags = logic.ToTagsMap(tagArr)
	}

	if m.Position != "" {
		arr := utils.SplitCutset(m.Position, ",，")
		if len(arr) == 2 {
			d.Position = logic.ToDmPointRpc(&types.Point{cast.ToFloat64(arr[0]), cast.ToFloat64(arr[1])})
		} else {
			return nil, errors.New("设备位置坐标格式错误")
		}
	}

	if m.Address != "" {
		d.Address = utils.ToRpcNullString(&m.Address)
	}

	return d, nil
}
