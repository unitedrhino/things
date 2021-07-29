package logic

import (
	"context"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/shared/utils"
	"gitee.com/godLei6/things/src/dmsvr/dm"
	"gitee.com/godLei6/things/src/webapi/internal/svc"
	"gitee.com/godLei6/things/src/webapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
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
	dis := make([]*types.DeviceInfo, 0, len(resp.Info))
	for _, v := range resp.Info {
		di := &types.DeviceInfo{
			ProductID:   v.ProductID,          //产品id 只读
			DeviceName:  v.DeviceName,         //设备名称 读写
			CreatedTime: v.CreatedTime,        //创建时间 只读
			Secret:      v.Secret,             //设备秘钥 只读
			FirstLogin:  v.FirstLogin,         //激活时间 只读
			LastLogin:   v.LastLogin,          //最后上线时间 只读
			Version:     v.Version.GetValue(), // 固件版本  读写
			LogLevel:    v.LogLevel,           // 日志级别:1)关闭 2)错误 3)告警 4)信息 5)调试  读写
			Cert:        v.Cert,               // 设备证书  只读
		}
		dis = append(dis, di)
	}
	return &types.GetDeviceInfoResp{
		Total: resp.Total,
		Info:  dis,
		Num:   int64(len(dis)),
	}, nil

}
