package devicegrouplogic

import (
	"context"
	"github.com/i-Things/things/shared/def"
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

/*
	检查当前分组嵌套层数是否超限，是返回true 否则返回false
*/
func (l *GroupInfoCreateLogic) CheckGroupLevel(groupID int64, level int64) (bool, error) {
	//采用递归方式，根据当前分组id和层数上限综合判断
	if groupID == 1 {
		if level <= def.DeviceGroupLevel && level > 1 {
			return false, nil
		}
		if level <= 1 {
			return true, nil
		}
	}

	resp, err := l.svcCtx.GroupInfo.FindOne(l.ctx, groupID)
	if err != nil {
		l.Errorf("%s.CheckGroupInfo msg=not find group id is %d\n", utils.FuncName(), groupID)
		return false, errors.Database.AddDetail(err)
	}

	return l.CheckGroupLevel(resp.ParentID, level-1)
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

	//判断当前分组parentid 层数是否达到指定层数，达到则不允许创建分组
	f, err := l.CheckGroupLevel(in.ParentID, def.DeviceGroupLevel)
	if err != nil {
		l.Errorf("%s.CheckGroupLevel in=%v\n", utils.FuncName(), in)
		return nil, errors.Database.AddDetail(err)
	}
	if f {
		l.Errorf("%s.CheckGroupInfo msg=group level is over %d \n", utils.FuncName(), def.DeviceGroupLevel)
		return nil, errors.OutRange.WithMsgf("子分组嵌套不能超过%d层", def.DeviceGroupLevel)
	}

	_, err = l.svcCtx.GroupInfo.Insert(l.ctx, &mysql.DmGroupInfo{
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
