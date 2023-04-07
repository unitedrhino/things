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
	"github.com/xuri/excelize/v2"
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

func (l *MultiImportLogic) MultiImport(req *types.DeviceMultiImportReq, csv *excelize.File, demoCnt int64) (resp *types.DeviceMultiImportResp, err error) {
	sm := sync.Map{}
	var egg errgroup.Group
	var errDatas []types.DeviceMultiImportErrdata

	rows, err := csv.Rows(csv.GetSheetName(0))
	if err != nil {
		return nil, errors.Parameter.WithMsg("读取表格Sheet失败了:" + err.Error())
	}

	rowCnt := int64(1)
	for ; rows.Next(); rowCnt++ {
		rowIdx := rowCnt
		if rowIdx <= demoCnt { //第1行是标题、第2、3行是示例，均跳过
			continue
		}

		//读取行数据
		cell, err := rows.Columns()
		if err != nil {
			errDatas = append(errDatas, types.DeviceMultiImportErrdata{rowIdx, "读取行数据出错:" + err.Error()})
			continue
		}
		//cell转dto
		dmReq, err := l.deviceMultiImportCellToDeviceInfo(cell)
		if err != nil {
			errDatas = append(errDatas, types.DeviceMultiImportErrdata{rowIdx, "行数据解析出错:" + err.Error()})
			continue
		}

		egg.Go(func() error {
			_, err := l.svcCtx.DeviceM.DeviceInfoCreate(l.ctx, dmReq)
			if err != nil {
				sm.Store(rowIdx, errors.Fmt(err).Msg)
			}
			return nil
		})

	} //end for
	rowCnt-- //最后多累计了1，要扣掉

	//阻塞等待所有gorouting
	if err = egg.Wait(); err != nil {
		return nil, err
	}

	sm.Range(func(key, value any) bool {
		errDatas = append(errDatas, types.DeviceMultiImportErrdata{key.(int64), "创建设备出错:" + value.(string)})
		return true
	})

	return &types.DeviceMultiImportResp{
		Total:   rowCnt - demoCnt,
		Errdata: errDatas,
	}, nil
}

//deviceMultiImportCellToDeviceInfo cell转dto
func (l *MultiImportLogic) deviceMultiImportCellToDeviceInfo(cell []string) (info *dm.DeviceInfo, err error) {
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
		return nil, errors.Parameter.WithMsg("缺少必填的产品名称")
	} else {
		d.ProductName = m.ProductName
	}

	if m.DeviceName == "" {
		return nil, errors.Parameter.WithMsg("缺少必填的设备名称")
	} else {
		d.DeviceName = m.DeviceName
	}

	if m.LogLevel != "" {
		if level, ok := def.LogLevelTextToIntMap[m.LogLevel]; !ok {
			return nil, errors.Parameter.WithMsg("设备日志级别格式错误")
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
				return nil, errors.Parameter.WithMsg("设备标签格式错误")
			}
		}
		d.Tags = logic.ToTagsMap(tagArr)
	}

	if m.Position != "" {
		arr := utils.SplitCutset(m.Position, ",，")
		if len(arr) == 2 {
			d.Position = logic.ToDmPointRpc(&types.Point{cast.ToFloat64(arr[0]), cast.ToFloat64(arr[1])})
		} else {
			return nil, errors.Parameter.WithMsg("设备位置坐标格式错误")
		}
	}

	if m.Address != "" {
		d.Address = utils.ToRpcNullString(&m.Address)
	}

	return d, nil
}
