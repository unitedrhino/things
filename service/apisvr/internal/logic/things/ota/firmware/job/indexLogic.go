package job

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.OtaFirmwareJobIndexReq) (resp *types.OtaFirmwareJobInfoIndexResp, err error) {
	var firmwareIndexReq dm.OtaFirmwareJobIndexReq
	_ = copier.Copy(&firmwareIndexReq, &req)
	index, err := l.svcCtx.OtaM.OtaFirmwareJobIndex(l.ctx, &firmwareIndexReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.OtaFirmwareIndex req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	var list []*types.OtaFirmwareJobInfo
	for _, v := range index.List {
		var result = types.OtaFirmwareJobInfo{}
		_ = copier.Copy(&result, &v)
		copier.Copy(&result.OtaFirmwareJobStatic, &v.Static)
		copier.Copy(&result.OtaFirmwareJobDynamic, &v.Dynamic)
		list = append(list, &result)
	}
	return &types.OtaFirmwareJobInfoIndexResp{List: list, Total: index.Total}, nil
}
