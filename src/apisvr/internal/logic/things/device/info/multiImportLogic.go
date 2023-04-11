package info

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/sync/errgroup"
	"strings"
	"sync"
)

type MultiImportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

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

func NewMultiImportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiImportLogic {
	return &MultiImportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiImportLogic) MultiImport(req *types.DeviceMultiImportReq, rows [][]string) (resp *types.DeviceMultiImportResp, err error) {
	sm := sync.Map{}
	var egg errgroup.Group
	var errDatas []types.DeviceMultiImportErrdata

	for i, cell := range rows {
		idx := i
		if idx == 0 {
			continue //第一行是 header，跳过
		}

		//cell转dto
		dmReq, err := l.deviceMultiImportCellToDeviceInfo(cell)
		if err != nil {
			errDatas = append(errDatas, types.DeviceMultiImportErrdata{int64(idx + 1), "行数据解析出错:" + errors.Fmt(err).GetMsg()})
			continue
		}

		egg.Go(func() error {
			_, err := l.svcCtx.DeviceM.DeviceInfoCreate(l.ctx, dmReq)
			if err != nil {
				sm.Store(idx, errors.Fmt(err).GetMsg())
			}
			return nil
		})

	} //end for

	//阻塞等待所有gorouting
	if err = egg.Wait(); err != nil {
		return nil, err
	}

	sm.Range(func(key, value any) bool {
		errDatas = append(errDatas, types.DeviceMultiImportErrdata{key.(int64), "创建设备出错:" + value.(string)})
		return true
	})

	return &types.DeviceMultiImportResp{
		Total:   int64(len(rows) - 1),
		Errdata: errDatas,
	}, nil
}

//deviceMultiImportCellToDeviceInfo cell转dto
func (l *MultiImportLogic) deviceMultiImportCellToDeviceInfo(cell []string) (info *dm.DeviceInfo, err error) {
	var (
		demoDataTag = "ithingsdemo"
		deviceInfo  = &dm.DeviceInfo{}
		importRow   = &multiImportCsvRow{
			ProductName: strings.TrimSpace(utils.SliceIndex(cell, 0, "")),
			DeviceName:  strings.TrimSpace(utils.SliceIndex(cell, 1, "")),
			LogLevel:    strings.TrimSpace(utils.SliceIndex(cell, 2, "")),
			Tags:        strings.TrimSpace(utils.SliceIndex(cell, 3, "")),
			Position:    strings.TrimSpace(utils.SliceIndex(cell, 4, "")),
			Address:     strings.TrimSpace(utils.SliceIndex(cell, 5, "")),
		}
	)

	if importRow.ProductName == "" {
		return nil, errors.Parameter.WithMsg("缺少必填的产品名称")
	} else {
		deviceInfo.ProductName = importRow.ProductName
	}

	if importRow.DeviceName == "" {
		return nil, errors.Parameter.WithMsg("缺少必填的设备名称")
	} else {
		deviceInfo.DeviceName = importRow.DeviceName
	}

	if strings.Contains(strings.ToLower(importRow.ProductName), demoDataTag) ||
		strings.Contains(strings.ToLower(importRow.DeviceName), demoDataTag) {
		return nil, errors.Parameter.WithMsg("请勿上传模板Demo数据")
	}

	if importRow.LogLevel != "" {
		if level, ok := def.LogLevelTextToIntMap[importRow.LogLevel]; !ok {
			return nil, errors.Parameter.WithMsg("设备日志级别格式错误")
		} else {
			deviceInfo.LogLevel = level
		}
	}

	if importRow.Tags != "" {
		arr := utils.SplitCutset(importRow.Tags, ";；\n")
		tagArr := make([]*types.Tag, 0, len(arr))
		for _, item := range arr {
			tagSli := utils.SplitCutset(item, ":：")
			if len(tagSli) == 2 {
				tagArr = append(tagArr, &types.Tag{tagSli[0], tagSli[1]})
			} else {
				return nil, errors.Parameter.WithMsg("设备标签格式错误")
			}
		}
		deviceInfo.Tags = logic.ToTagsMap(tagArr)
	}

	if importRow.Position != "" {
		arr := utils.SplitCutset(importRow.Position, ",，")
		if len(arr) == 2 {
			deviceInfo.Position = logic.ToDmPointRpc(&types.Point{cast.ToFloat64(arr[0]), cast.ToFloat64(arr[1])})
		} else {
			return nil, errors.Parameter.WithMsg("设备位置坐标格式错误")
		}
	}

	if importRow.Address != "" {
		deviceInfo.Address = utils.ToRpcNullString(&importRow.Address)
	}

	return deviceInfo, nil
}
