package logic

import (
	"context"
	"gitee.com/godLei6/things/shared/def"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/src/dcsvr/dc"
	"gitee.com/godLei6/things/src/dcsvr/internal/svc"
	"gitee.com/godLei6/things/src/dcsvr/model"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMemberLogic {
	return &GetGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取组成员
func (l *GetGroupMemberLogic) GetGroupMember(in *dc.GetGroupMemberReq) (resp *dc.GetGroupMemberResp,err error) {
	l.Infof("GetGroupMember|req=%+v", in)
	var info []*dc.GroupMember
	var size int64
	if (in.Page == nil || in.Page.Page == 0) && (in.GroupID == 0 && in.MemberID == ""){
		di, err := l.svcCtx.GroupMember.FindOneByGroupIDMemberID(in.GroupID, in.MemberID)
		if err != nil {
			if err == model.ErrNotFound {
				return nil, errors.NotFind
			}
			return nil, err
		}
		info = append(info, DBToRPCFmt(di).(*dc.GroupMember))
	}else {
		var page int64 = 1
		var pageSize int64 = 20
		var di []*model.GroupMember
		if !(in.Page == nil || in.Page.Page == 0 || in.Page.PageSize == 0) {
			page = in.Page.Page
			pageSize = in.Page.PageSize
		}
		if in.GroupID != 0 {
			size, err = l.svcCtx.DcDB.GetCountByGroupMemberGroupID(
				in.GroupID)
			if err != nil {
				return nil, err
			}
			di, err = l.svcCtx.DcDB.FindByGroupMemberGroupID(
				in.GroupID, def.PageInfo{PageSize: pageSize, Page: page})
			if err != nil {
				return nil, err
			}
		}else{
			size, err = l.svcCtx.DcDB.GetCountByGroupMemberMemberID(
				in.MemberID)
			if err != nil {
				return nil, err
			}
			di, err = l.svcCtx.DcDB.FindByGroupMemberMemberID(
				in.MemberID, def.PageInfo{PageSize: pageSize, Page: page})
			if err != nil {
				return nil, err
			}
		}
		info = make([]*dc.GroupMember, 0, len(di))
		for _, v := range di {
			info = append(info, DBToRPCFmt(v).(*dc.GroupMember))
		}
	}
	return &dc.GetGroupMemberResp{Info: info, Total: size}, nil
}
