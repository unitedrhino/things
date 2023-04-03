package info

import (
	"bytes"
	"context"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
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
		if cell, err := rows.Columns(); err != nil {
			errDatas = append(errDatas, types.DeviceMultiImportErrdata{
				Row: r,
				Msg: "读取行数据出错:" + err.Error(),
			})
		} else {
			if r < 2 { //第1行是标题、第2行是示例，均跳过
				continue
			}

			//TODO 校验行单元格数据

			egg.Go(func() error {
				_ = cell //TODO 异步添加设备逻辑
				return nil
			})
		}
	}
	//阻塞等待所有gorouting
	if err = egg.Wait(); err != nil {
		return nil, err
	}

	return
}
