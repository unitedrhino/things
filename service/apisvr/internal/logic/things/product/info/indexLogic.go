package info

import (
	"context"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/apisvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

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

func (l *IndexLogic) Index(req *types.ProductInfoIndexReq) (resp *types.ProductInfoIndexResp, err error) {
	dmReq := &dm.ProductInfoIndexReq{
		DeviceType:   req.DeviceType, //产品id
		ProductName:  req.ProductName,
		ProductIDs:   req.ProductIDs,
		Tags:         logic.ToTagsMap(req.Tags),
		Page:         logic.ToDmPageRpc(req.Page),
		ProtocolCode: req.ProtocolCode,
		WithProtocol: req.WithProtocol,
		ProtocolType: req.ProtocolType,
		WithCategory: req.WithCategory,
		DeviceTypes:  req.DeviceTypes,
		SceneMode:    req.SceneMode,
		Statuses:     req.Statuses,
		Status:       req.Status,
		NetType:      req.NetType,
		SceneModes:   req.SceneModes,
		ProjectID:    req.ProjectID,
		AreaID:       req.AreaID,
		AreaIDPath:   req.AreaIDPath,
		CategoryIDs:  req.CategoryIDs,
	}
	dmResp, err := l.svcCtx.ProductM.ProductInfoIndex(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.GetDeviceInfo req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	pis := make([]*types.ProductInfo, 0, len(dmResp.List))
	for _, v := range dmResp.List {
		pi := productInfoToApi(l.ctx, v)
		pis = append(pis, pi)
	}
	return &types.ProductInfoIndexResp{
		PageResp: logic.ToPageResp(req.Page, dmResp.Total),
		List:     pis,
	}, nil
}
