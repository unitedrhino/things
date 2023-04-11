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
	"sort"
	"strings"
	"sync"
)

type MultiImportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMultiImportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiImportLogic {
	return &MultiImportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiImportLogic) MultiImport(req *types.DeviceMultiImportReq, rows [][]string) (resp *types.DeviceMultiImportResp, err error) {
	var (
		sm      = sync.Map{}
		egg     errgroup.Group
		headers *types.DeviceMultiImportRow
		errdata []*types.DeviceMultiImportRow
	)

	for i, cell := range rows {
		idx := int64(i) //这里必须是 int64，因为下面 key.(int64) 要推断

		//cell转dto
		importRow := l.deviceMultiImportRowToDto(idx, cell)
		if idx == 0 { //第一行是表头
			headers = importRow
			continue //数据处理 跳过表头
		}

		dmDeviceInfoReq, err := l.deviceMultiImportRowToDeviceInfo(importRow)
		if err != nil {
			importRow.Tips = "行数据解析出错:" + errors.Fmt(err).GetMsg()
			errdata = append(errdata, importRow)
			continue
		}

		egg.Go(func() error {
			_, err := l.svcCtx.DeviceM.DeviceInfoCreate(l.ctx, dmDeviceInfoReq)
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

	sm.Range(func(i, value any) bool {
		idx := i.(int64) //这里必须是 int64，因为下面 key.(int64) 要推断
		importRow := l.deviceMultiImportRowToDto(idx, rows[idx])
		importRow.Tips = "创建设备出错:" + value.(string)
		errdata = append(errdata, importRow)
		return true
	})

	if len(errdata) > 0 { //重新排序
		sort.Slice(errdata, func(i, j int) bool { return errdata[i].Row < errdata[j].Row })
	}

	return &types.DeviceMultiImportResp{
		Total:   int64(len(rows) - 1),
		Headers: headers,
		Errdata: errdata,
	}, nil
}

func (l *MultiImportLogic) deviceMultiImportRowToDto(idx int64, cell []string) *types.DeviceMultiImportRow {
	return &types.DeviceMultiImportRow{
		Row:         idx,
		ProductName: strings.TrimSpace(utils.SliceIndex(cell, 0, "")),
		DeviceName:  strings.TrimSpace(utils.SliceIndex(cell, 1, "")),
		LogLevel:    strings.TrimSpace(utils.SliceIndex(cell, 2, "")),
		Tags:        strings.TrimSpace(utils.SliceIndex(cell, 3, "")),
		Position:    strings.TrimSpace(utils.SliceIndex(cell, 4, "")),
		Address:     strings.TrimSpace(utils.SliceIndex(cell, 5, "")),
		Tips:        strings.TrimSpace(utils.SliceIndex(cell, 6, "")),
	}
}

//deviceMultiImportRowToDeviceInfo cell转dto
func (l *MultiImportLogic) deviceMultiImportRowToDeviceInfo(importRow *types.DeviceMultiImportRow) (info *dm.DeviceInfo, err error) {
	var (
		demoDataTag = "ithingsdemo"
		deviceInfo  = &dm.DeviceInfo{}
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
