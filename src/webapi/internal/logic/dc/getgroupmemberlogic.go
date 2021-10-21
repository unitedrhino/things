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

type GetGroupMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetGroupMemberLogic {
	return GetGroupMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGroupMemberLogic) GetGroupMember(req types.GetGroupMemberReq) (*types.GetGroupMemberResp, error) {
	l.Infof("GetGroupMember|req=%+v", req)
	dcReq := &dc.GetGroupMemberReq{
		GroupID: req.GroupID,
		MemberID:req.MemberID,
		MemberType:req.MemberType,
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
	resp, err := l.svcCtx.DcRpc.GetGroupMember(l.ctx, dcReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.GetGroupMember|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	gis := make([]*types.GroupMember, 0, len(resp.Info))
	for _, v := range resp.Info {
		gi := RPCToApiFmt(v).(*types.GroupMember)
		gis = append(gis, gi)
	}
	return &types.GetGroupMemberResp{
		Total: resp.Total,
		Info:  gis,
		Num:   int64(len(gis)),
	}, nil
}
