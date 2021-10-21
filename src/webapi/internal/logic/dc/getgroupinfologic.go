package dc

import (
	"context"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/shared/utils"
	"gitee.com/godLei6/things/src/dcsvr/dc"
	"gitee.com/godLei6/things/src/webapi/internal/svc"
	"gitee.com/godLei6/things/src/webapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetGroupInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetGroupInfoLogic {
	return GetGroupInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGroupInfoLogic) GetGroupInfo(req types.GetGroupInfoReq) (*types.GetGroupInfoResp, error) {
	l.Infof("GetGroupInfo|req=%+v", req)
	dcReq := &dc.GetGroupInfoReq{
		GroupID: req.GroupID,
	}
	if req.Page != nil {
		if req.Page.PageSize == 0 || req.Page.Page == 0 {
			return nil, errors.Parameter.AddDetail("pageSize and page can't equal 0")
		}
		dcReq.Page = &dc.PageInfo{
			Page:     req.Page.Page,
			PageSize: req.Page.PageSize,
		}
	}
	resp, err := l.svcCtx.DcRpc.GetGroupInfo(l.ctx, dcReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.GetGroupInfo|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	gis := make([]*types.GroupInfo, 0, len(resp.Info))
	for _, v := range resp.Info {
		gi := RPCToApiFmt(v).(*types.GroupInfo)
		gis = append(gis, gi)
	}
	return &types.GetGroupInfoResp{
		Total: resp.Total,
		Info:  gis,
		Num:   int64(len(gis)),
	}, nil
}
