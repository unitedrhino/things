package data

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/dm"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SchemaLatestIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSchemaLatestIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SchemaLatestIndexLogic {
	return &SchemaLatestIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SchemaLatestIndexLogic) SchemaLatestIndex(req *types.DataSchemaLatestIndexReq) (resp *types.DataSchemaIndexResp, err error) {
	dmResp, err := l.svcCtx.DmRpc.DataSchemaLatestIndex(l.ctx, &dm.DataSchemaLatestIndexReq{
		Method:     req.Method,
		DeviceName: req.DeviceName,
		ProductID:  req.ProductID,
		DataID:     req.DataID,
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
