package rulelogic

import (
	"context"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"
	"gorm.io/gorm"

	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmInfoCreateLogic {
	return &AlarmInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AlarmInfoCreateLogic) AlarmInfoCreate(in *ud.AlarmInfo) (*ud.WithID, error) {
	po := utils.Copy[relationDB.UdAlarmInfo](in)
	po.ID = 0
	err := stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		err := relationDB.NewAlarmInfoRepo(l.ctx).Insert(l.ctx, po)
		if err != nil {
			return err
		}
		if len(in.SceneIDs) != 0 {
			var pos []*relationDB.UdAlarmScene
			for _, v := range in.SceneIDs {
				pos = append(pos, &relationDB.UdAlarmScene{SceneID: v, AlarmID: po.ID})
			}
			err = relationDB.NewAlarmSceneRepo(l.ctx).MultiInsert(l.ctx, pos)
		}
		return err
	})

	return &ud.WithID{Id: po.ID}, err
}
