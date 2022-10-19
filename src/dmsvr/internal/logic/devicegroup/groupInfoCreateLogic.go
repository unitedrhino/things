package devicegrouplogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupInfoCreateLogic {
	return &GroupInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

/*
发现返回true 没有返回false
*/
func (l *GroupInfoCreateLogic) CheckGroupInfo(in *dm.GroupInfoCreateReq) (bool, error) {
	_, err := l.svcCtx.GroupInfo.FindOneByGroupName(l.ctx, in.GroupName)
	switch err {
	case mysql.ErrNotFound:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}
}

// 创建分组
func (l *GroupInfoCreateLogic) GroupInfoCreate(in *dm.GroupInfoCreateReq) (*dm.Response, error) {
	find, err := l.CheckGroupInfo(in)
	if err != nil {
		l.Errorf("%s.CheckGroupInfo in=%v\n", utils.FuncName(), in)
		return nil, errors.Database.AddDetail(err)
	} else if find == true {
		return nil, errors.Duplicate.WithMsgf("组名重复:%s", in.GroupName).AddDetail("GroupName:" + in.GroupName)
	}
	_, err = l.svcCtx.GroupInfo.Insert(l.ctx, &mysql.GroupInfo{
		GroupID:   l.svcCtx.GroupID.GetSnowflakeId(),
		ParentID:  in.ParentID,
		GroupName: in.GroupName,
		Desc:      in.Desc,
		Tags:      "{}",
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}

	return &dm.Response{}, nil
}
