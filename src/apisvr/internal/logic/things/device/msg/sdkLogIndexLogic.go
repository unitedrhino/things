package msg

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/disvr/pb/di"

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

func (l *SdkLogIndexLogic) SdkLogIndex(req *types.DeviceMsgSdkLogIndexReq) (resp *types.DeviceMsgSdkIndexResp, err error) {
	dmReq := &di.SdkLogIndexReq{
		DeviceName: req.DeviceName,
		ProductID:  req.ProductID, //产品id
		TimeStart:  req.TimeStart,
		TimeEnd:    req.TimeEnd,
		Page: &di.PageInfo{
			Page: req.Page.Page,
			Size: req.Page.Size,
		},
	}

	dmResp, err := l.svcCtx.DeviceMsg.SdkLogIndex(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.GetDeviceSDKLog|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	info := make([]*types.DeviceMsgSdkIndex, 0, len(dmResp.List))
	for _, v := range dmResp.List {
		info = append(info, &types.DeviceMsgSdkIndex{
			Timestamp: v.Timestamp,
			Loglevel:  v.Loglevel,
			Content:   v.Content,
		})
	}
	return &types.DeviceMsgSdkIndexResp{List: info, Total: dmResp.Total}, err

}
