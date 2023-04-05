package info

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/apisvr/internal/logic/things/device"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/xuri/excelize/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/sync/errgroup"
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

const MultiImportDemoCnt = int64(3)

func (l *MultiImportLogic) MultiImport(req *types.DeviceMultiImportReq, csv *excelize.File, demoCnt int64) (resp *types.DeviceMultiImportResp, err error) {
	sm := sync.Map{}
	var egg errgroup.Group
	var errDatas []types.DeviceMultiImportErrdata

	rows, err := csv.Rows(csv.GetSheetName(0))
	if err != nil {
		return nil, errors.Parameter.WithMsg("读取表格Sheet失败了:" + err.Error())
	}

	var rowCnt int64
	for ; rows.Next(); rowCnt++ {
		rowIdx := rowCnt
		if rowIdx < demoCnt { //第1行是标题、第2、3行是示例，均跳过
			continue
		}

		//读取行数据
		cell, err := rows.Columns()
		if err != nil {
			errDatas = append(errDatas, types.DeviceMultiImportErrdata{rowIdx, "读取行数据出错:" + err.Error()})
			continue
		}
		//cell转dto
		dmReq, err := device.DeviceMultiImportCellToDeviceInfo(cell)
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
