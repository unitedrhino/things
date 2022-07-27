package logic

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dcsvr/internal/repo/mysql"
	"time"

	"github.com/i-Things/things/src/dcsvr/dc"
	"github.com/i-Things/things/src/dcsvr/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ManageGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewManageGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ManageGroupMemberLogic {
	return &ManageGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

/*
发现返回true 没有返回false
*/
func (l *ManageGroupMemberLogic) CheckGroup(in *dc.ManageGroupMemberReq) (bool, error) {
	_, err := l.svcCtx.GroupInfo.FindOne(in.Info.GroupID)
	switch err {
	case mysql.ErrNotFound:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}
}

/*
发现返回true 没有返回false
*/
func (l *ManageGroupMemberLogic) CheckGroupMember(in *dc.ManageGroupMemberReq) (bool, error) {
	_, err := l.svcCtx.GroupMember.FindOneByGroupIDMemberIDMemberType(
		in.Info.GroupID, in.Info.MemberID, in.Info.MemberType)
	switch err {
	case mysql.ErrNotFound:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}
}

func (l *ManageGroupMemberLogic) AddGroupMember(in *dc.ManageGroupMemberReq) (*dc.GroupMember, error) {
	find, err := l.CheckGroupMember(in)
	if err != nil {
		return nil, errors.System.AddDetail(err)
	} else if find == true {
		return nil, errors.Duplicate.AddDetailf(
			"GroupID:%v, MemberID:%v,MemberType:%v", in.Info.GroupID, in.Info.MemberID, in.Info.MemberType)
	}
	l.Infof("find=%v|err=%v\n", find, err)
	find, err = l.CheckGroup(in)
	if err != nil {
		return nil, errors.System.AddDetail(err)
	} else if find == false {
		return nil, errors.Parameter.AddDetail(
			"not find GroupID:%v, MemberID:%v,MemberType:%v",
			in.Info.GroupID, in.Info.MemberID, in.Info.MemberType)
	}

	di := mysql.GroupMember{
		GroupID:     in.Info.GroupID,    // 组id
		MemberID:    in.Info.MemberID,   // 成员id
		MemberType:  in.Info.MemberType, // 成员类型:1:设备 2:用户
		CreatedTime: time.Now(),
	}
	if in.Info.MemberType > 2 || in.Info.MemberType < 1 {
		return nil, errors.Parameter.AddDetail(
			"MemberType not support:", in.Info.MemberType)
	}
	_, err = l.svcCtx.GroupMember.Insert(di)
	if err != nil {
		l.Errorf("AddDevice|DeviceInfo|Insert|err=%+v", err)
		return nil, errors.System.AddDetail(err)
	}
	return DBToRPCFmt(&di).(*dc.GroupMember), nil
}

func (l *ManageGroupMemberLogic) DelGroupMember(in *dc.ManageGroupMemberReq) (*dc.GroupMember, error) {
	di, err := l.svcCtx.GroupMember.FindOneByGroupIDMemberIDMemberType(
		in.Info.GroupID, in.Info.MemberID, in.Info.MemberType)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.Parameter.AddDetail(
				fmt.Sprintf("not find GroupMember|GroupID=%v|MemberID=%sMemberType=%d",
					in.Info.GroupID, in.Info.MemberID, in.Info.MemberType))
		}
		l.Errorf("DelGroupMember|GroupMember|FindOne|err=%+v", err)
		return nil, errors.System.AddDetail(err)
	}
	err = l.svcCtx.GroupMember.Delete(di.Id)
	if err != nil {
		l.Errorf("DelGroupMember|GroupMember|Delete|err=%+v", err)
		return nil, errors.System.AddDetail(err)
	}
	return &dc.GroupMember{}, nil
}

// 管理组成员
func (l *ManageGroupMemberLogic) ManageGroupMember(in *dc.ManageGroupMemberReq) (*dc.GroupMember, error) {
	defer func() {
		if p := recover(); p != nil {
			utils.HandleThrow(l.ctx, p)
		}
	}()
	l.Infof("ManageGroupMember|req=%+v", in)
	switch in.Opt {
	case def.OPT_ADD:
		return l.AddGroupMember(in)
	case def.OPT_DEL:
		return l.DelGroupMember(in)
	default:
		return nil, errors.Parameter.AddDetail("not support opt:" + string(in.Opt))
	}

	return &dc.GroupMember{}, nil
}
