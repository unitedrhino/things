package data

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SchemaLogIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSchemaLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SchemaLogIndexLogic {
	return &SchemaLogIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SchemaLogIndexLogic) SchemaLogIndex(req *types.DataSchemaLogIndexReq) (resp *types.DataSchemaIndexResp, err error) {
	dmResp, err := l.svcCtx.DmRpc.DataSchemaLogIndex(l.ctx, &dm.DataSchemaLogIndexReq{
		Method:     req.Method,
		DeviceName: req.DeviceName,
		ProductID:  req.ProductID,
		DataID:     req.DataID,
		TimeStart:  req.TimeStart,
		TimeEnd:    req.TimeEnd,
		Page: &dm.PageInfo{
			Page: req.Page.Page,
			Size: req.Page.Size,
		},
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.GetDeviceData|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	info := make([]*types.DataSchemaIndex, 0, len(dmResp.List))
	for _, v := range dmResp.List {
		info = append(info, &types.DataSchemaIndex{
			Timestamp: v.Timestamp,
			Type:      v.Type,
			DataID:    v.DataID,
			GetValue:  v.GetValue,
			SendValue: v.SendValue,
		})
	}
	return &types.DataSchemaIndexResp{
		Total: dmResp.Total,
		List:  info,
	}, nil
}
