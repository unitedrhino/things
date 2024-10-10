package rulelogic

import (
	"context"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/udsvr/internal/repo/relationDB"
	"gorm.io/gorm"

	"gitee.com/unitedrhino/things/service/udsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmInfoUpdateLogic {
	return &AlarmInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AlarmInfoUpdateLogic) AlarmInfoUpdate(in *ud.AlarmInfo) (*ud.Empty, error) {
	old, err := relationDB.NewAlarmInfoRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if in.Name != "" {
		old.Name = in.Name
	}
	if in.Desc != "" {
		old.Desc = in.Desc
	}
	if in.Level != 0 {
		old.Level = in.Level
	}
	if in.Status != 0 {
		old.Status = in.Status
	}
	if in.UserIDs != nil {
		old.UserIDs = in.UserIDs
	}
	if len(in.Notifies) != 0 {
		old.Notifies = utils.CopySlice[relationDB.UdAlarmNotify](in.Notifies)
	}
	err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		err = relationDB.NewAlarmInfoRepo(tx).Update(l.ctx, old)
		if err != nil {
			return err
		}
		if in.SceneIDs != nil {
			asDB := relationDB.NewAlarmSceneRepo(tx)
			err := asDB.DeleteByFilter(l.ctx, relationDB.AlarmSceneFilter{
				AlarmID: old.ID,
			})
			if err != nil {
				return err
			}
			if len(in.SceneIDs) > 0 {
				var pos []*relationDB.UdAlarmScene
				for _, v := range in.SceneIDs {
					pos = append(pos, &relationDB.UdAlarmScene{SceneID: v, AlarmID: old.ID})
				}
				err = asDB.MultiInsert(l.ctx, pos)
			}
			return err
		}
		return nil
	})
	return &ud.Empty{}, err
}
