package info

import (
	"bytes"
	"context"
	"github.com/i-Things/things/src/apisvr/internal/logic/things/device"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/xuri/excelize/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/sync/errgroup"
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

func (l *MultiImportLogic) MultiImport(req *types.DeviceMultiImportReq) (resp *types.DeviceMultiImportResp, err error) {
	reader := bytes.NewReader(req.File)
	csv, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, err
	}
	rows, err := csv.Rows("Sheet1")
	if err != nil {
		return nil, err
	}

	var egg errgroup.Group
	var errDatas []types.DeviceMultiImportErrdata

	for r := int64(0); rows.Next(); r++ {
		if r < 3 { //第1行是标题、第2、3行是示例，均跳过
			continue
		}

		//读取行数据
		cell, err := rows.Columns()
		if err != nil {
			errDatas = append(errDatas, types.DeviceMultiImportErrdata{
				Row: r,
				Msg: "读取行数据出错:" + err.Error(),
			})
			return nil, err
		}

		//解析行数据到Dto，并验证数据格式
		rowDto := device.NewMultiImportCsvRow(cell)
		if err := rowDto.Valid(); err != nil {
			return nil, err
		}

		egg.Go(func() error {
			dmReq := &dm.DeviceInfo{
				//ProductID:  req.ProductID,  //产品id 只读
				//DeviceName: req.DeviceName, //设备名称 读写
				//LogLevel:   req.LogLevel,   // 日志级别:1)关闭 2)错误 3)告警 4)信息 5)调试  读写
				//Tags:       logic.ToTagsMap(req.Tags),
				//Address:    utils.ToRpcNullString(req.Address),
				//Position:   logic.ToDmPointRpc(req.Position),
			}
			_, err := l.svcCtx.DeviceM.DeviceInfoCreate(l.ctx, dmReq)
			if err != nil {

			}
			return nil
		})
	}
	//阻塞等待所有gorouting
	if err = egg.Wait(); err != nil {
		return nil, err
	}

	return
}
