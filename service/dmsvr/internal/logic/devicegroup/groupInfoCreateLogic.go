package devicegrouplogic

import (
	"context"
	"gitee.com/i-Things/core/shared/def"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	GiDB *relationDB.GroupInfoRepo
}

func NewGroupInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupInfoCreateLogic {
	return &GroupInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		GiDB:   relationDB.NewGroupInfoRepo(ctx),
	}
}

/*
发现返回true 没有返回false
*/
func (l *GroupInfoCreateLogic) CheckGroupInfo(in *dm.GroupInfo) (bool, error) {
	_, err := l.GiDB.FindOneByFilter(l.ctx, relationDB.GroupInfoFilter{Names: []string{in.Name}})
	if err == nil {
		return true, nil
	}
	if errors.Cmp(err, errors.NotFind) {
		return false, nil
	}
	return false, err
}

/*
检查当前分组嵌套层数是否超限，是返回true 否则返回false
*/
func (l *GroupInfoCreateLogic) CheckGroupLevel(ID int64, level int64) (bool, error) {
	//采用递归方式，根据当前分组id和层数上限综合判断
	if ID == 1 {
		if level <= def.DeviceGroupLevel && level > 1 {
			return false, nil
		}
		if level <= 1 {
			return true, nil
		}
	}

	resp, err := l.GiDB.FindOne(l.ctx, ID)
	if err != nil {
		l.Errorf("%s.CheckGroupInfo msg=not find group id is %d\n", utils.FuncName(), ID)
		return false, errors.Database.AddDetail(err)
	}

	return l.CheckGroupLevel(resp.ParentID, level-1)
}

// 创建分组
func (l *GroupInfoCreateLogic) GroupInfoCreate(in *dm.GroupInfo) (*dm.WithID, error) {
	if in.AreaID == 0 {
		in.AreaID = def.NotClassified
	}
	find, err := l.CheckGroupInfo(in)
	if err != nil {
		l.Errorf("%s.CheckGroupInfo in=%v\n", utils.FuncName(), in)
		return nil, errors.Database.AddDetail(err)
	} else if find == true {
		return nil, errors.Duplicate.WithMsgf("组名重复:%s", in.Name).AddDetail("Name:" + in.Name)
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
	po := relationDB.DmGroupInfo{
		ParentID:  in.ParentID,
		ProductID: in.ProductID,
		Name:      in.Name,
		Desc:      in.Desc,
		Tags:      in.Tags,
	}
	err = l.GiDB.Insert(l.ctx, &po)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}

	return &dm.WithID{Id: po.ID}, nil
}
