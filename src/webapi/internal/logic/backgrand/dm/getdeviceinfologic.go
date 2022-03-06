package dm

import (
	"context"
	"github.com/go-things/things/shared/errors"
	"github.com/go-things/things/shared/utils"
	"github.com/go-things/things/src/dmsvr/dm"
	"github.com/go-things/things/src/webapi/internal/svc"
	"github.com/go-things/things/src/webapi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDeviceInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDeviceInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetDeviceInfoLogic {
	return GetDeviceInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDeviceInfoLogic) GetDeviceInfo(req types.GetDeviceInfoReq) (*types.GetDeviceInfoResp, error) {
	l.Infof("GetDeviceInfo|req=%+v", req)
	dmReq := &dm.GetDeviceInfoReq{
		DeviceName: req.DeviceName, //设备名 为空时获取产品id下的所有设备信息
		ProductID:  req.ProductID,  //产品id
	}
	if req.Page != nil {
		if req.Page.PageSize == 0 || req.Page.Page == 0 {
			return nil, errors.Parameter.AddDetail("pageSize and page can't equal 0")
		}
		dmReq.Page = &dm.PageInfo{
			Page:     req.Page.Page,
			PageSize: req.Page.PageSize,
		}
	}
	resp, err := l.svcCtx.DmRpc.GetDeviceInfo(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.GetDeviceInfo|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	dis := make([]*types.DeviceInfo, 0, len(resp.List))
	for _, v := range resp.List {
		di := types.DeviceInfoToApi(v)
		dis = append(dis, di)
	}
	return &types.GetDeviceInfoResp{
		Total: resp.Total,
		List:  dis,
		Num:   int64(len(dis)),
	}, nil
}
