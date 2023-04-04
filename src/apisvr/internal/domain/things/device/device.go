package device

import (
	"errors"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/utils"
	"github.com/spf13/cast"
)

type (
	MultiImportCsvRow struct {
		ProductName string        //【必填】产品名称
		DeviceName  string        //【必填】设备名称 读写
		LogLevel    string        //【可选】日志级别:关闭 错误 告警 信息 5调试
		Tags        string        //【可选】设备tags
		Position    string        //【可选】设备定位,默认百度坐标系
		Address     string        //【可选】所在详细地址
		dtoLogLevel def.LogLevel  //日志级别:1)关闭 2)错误 3)告警 4)信息 5)调试
		dtoTags     []devices.Tag //设备tags，如 {"color":"黑色","fsef":"gfadfa","gadsfa":"agaef"}
		dtoPoint    devices.Point //设备位置坐标，如 116.442501,40.03446
	}
)

func (m MultiImportCsvRow) Valid() error {
	if m.ProductName == "" {
		return errors.New("缺少必填的产品名称")
	}
	if m.DeviceName == "" {
		return errors.New("缺少必填的设备名称")
	}

	if m.LogLevel != "" {
		if level, ok := def.LogLevelTextToIntMap[m.LogLevel]; !ok {
			return errors.New("设备日志级别格式错误")
		} else {
			m.dtoLogLevel = level
		}
	}

	if m.Tags != "" {
		arr := utils.SplitCutset(m.Tags, ";；\n")
		for _, item := range arr {
			tagSli := utils.SplitCutset(item, ":：")
			if len(tagSli) == 2 {
				m.dtoTags = append(m.dtoTags, devices.Tag{tagSli[0], tagSli[1]})
			} else {
				return errors.New("设备标签格式错误")
			}
		}
	}

	if m.Position != "" {
		arr := utils.SplitCutset(m.Position, ":：")
		if len(arr) == 2 {
			m.dtoPoint = devices.Point{cast.ToFloat64(arr[0]), cast.ToFloat64(arr[1])}
		} else {
			return errors.New("设备位置坐标格式错误")
		}
	}

	if m.Address != "" {
		//do nothing
	}

	return nil
}
