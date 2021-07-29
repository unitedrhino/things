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

type GetProductInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetProductInfoLogic {
	return GetProductInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProductInfoLogic) GetProductInfo(req types.GetProductInfoReq) (*types.GetProductInfoResp, error) {
	dmReq := &dm.GetProductInfoReq{
		ProductID: req.ProductID, //产品id
	}
	if req.Page != nil {
		if req.Page.PageSize == 0 || req.Page.Page == 0 {
			return nil, errors.Parameter.AddDetail("pageSize and page can't equal 0")
		}
		dmReq.Page = &dm.PageInfo{
			Page:     req.Page.Page,
			PageSize: req.Page.PageSize,
		}
	} else if req.ProductID == "" {
		return nil, errors.Parameter.AddDetail("need page or product")
	}
	resp, err := l.svcCtx.DmRpc.GetProductInfo(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.GetDeviceInfo|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	pis := make([]*types.ProductInfo, 0, len(resp.Info))
	for _, v := range resp.Info {
		pi := &types.ProductInfo{
			CreatedTime:  v.CreatedTime,            //创建时间 只读
			ProductID:    v.ProductID,              //产品id 只读
			ProductName:  v.ProductName,            //产品名称
			AuthMode:     v.AuthMode,               //认证方式:0:账密认证,1:秘钥认证
			DeviceType:   v.DeviceType,             //设备类型:0:设备,1:网关,2:子设备
			CategoryID:   v.CategoryID,             //产品品类
			NetType:      v.NetType,                //通讯方式:0:其他,1:wi-fi,2:2G/3G/4G,3:5G,4:BLE,5:LoRaWAN
			DataProto:    v.DataProto,              //数据协议:0:自定义,1:数据模板
			AutoRegister: v.AutoRegister,           //动态注册:0:关闭,1:打开,2:打开并自动创建设备
			Secret:       v.Secret,                 //动态注册产品秘钥 只读
			Template:     v.Template.GetValue(),    //数据模板
			Description:  v.Description.GetValue(), //描述
			DevStatus:    v.DevStatus.GetValue(),   // 产品状态
		}
		pis = append(pis, pi)
	}
	return &types.GetProductInfoResp{
		Total: resp.Total,
		Info:  pis,
		Num:   int64(len(pis)),
	}, nil
	return &types.GetProductInfoResp{}, nil
}
