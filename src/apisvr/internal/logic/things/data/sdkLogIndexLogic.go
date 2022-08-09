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

type SdkLogIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSdkLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkLogIndexLogic {
	return &SdkLogIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SdkLogIndexLogic) SdkLogIndex(req *types.DataSdkLogIndexReq) (resp *types.DataSdkIndexResp, err error) {
	dmReq := &dm.DataSdkLogIndexReq{
		DeviceName: req.DeviceName,
		ProductID:  req.ProductID, //产品id
		TimeStart:  req.TimeStart,
		TimeEnd:    req.TimeEnd,
		Page: &dm.PageInfo{
			Page: req.Page.Page,
			Size: req.Page.Size,
		},
	}

	dmResp, err := l.svcCtx.DmRpc.DataSdkLogIndex(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.GetDeviceSDKLog|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	info := make([]*types.DataSdkIndex, 0, len(dmResp.List))
	for _, v := range dmResp.List {
		info = append(info, &types.DataSdkIndex{
			Timestamp: v.Timestamp,
			Loglevel:  v.Loglevel,
			Content:   v.Content,
		})
	}
	return &types.DataSdkIndexResp{List: info, Total: dmResp.Total}, err
}
