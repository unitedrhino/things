package dm

import (
	"context"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/shared/utils"
	"gitee.com/godLei6/things/src/dmsvr/dm"
	"github.com/golang/protobuf/ptypes/wrappers"

	"gitee.com/godLei6/things/src/webapi/internal/svc"
	"gitee.com/godLei6/things/src/webapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ManageProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewManageProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) ManageProductLogic {
	return ManageProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ManageProductLogic) ManageProduct(req types.ManageProductReq) (*types.ProductInfo, error) {
	l.Infof("ManageProduct|req=%+v", req)
	dmReq := &dm.ManageProductReq{
		Opt: req.Opt,
		Info: &dm.ProductInfo{
			ProductID:    req.Info.ProductID,    //产品id 只读
			ProductName:  req.Info.ProductName,  //产品名称
			AuthMode:     req.Info.AuthMode,     //认证方式:0:账密认证,1:秘钥认证
			DeviceType:   req.Info.DeviceType,   //设备类型:0:设备,1:网关,2:子设备
			CategoryID:   req.Info.CategoryID,   //产品品类
			NetType:      req.Info.NetType,      //通讯方式:0:其他,1:wi-fi,2:2G/3G/4G,3:5G,4:BLE,5:LoRaWAN
			DataProto:    req.Info.DataProto,    //数据协议:0:自定义,1:数据模板
			AutoRegister: req.Info.AutoRegister, //动态注册:0:关闭,1:打开,2:打开并自动创建设备
		},
	}
	if req.Info.Template != nil {
		dmReq.Info.Template = &wrappers.StringValue{
			Value: *req.Info.Template,
		}
	}
	if req.Info.Description != nil {
		dmReq.Info.Description = &wrappers.StringValue{
			Value: *req.Info.Description,
		}
	}
	if req.Info.DevStatus != nil {
		dmReq.Info.DevStatus = &wrappers.StringValue{
			Value: *req.Info.DevStatus,
		}
	}
	resp, err := l.svcCtx.DmRpc.ManageProduct(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.ManageProduct|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return RPCToApiFmt(resp).(*types.ProductInfo), nil
}
