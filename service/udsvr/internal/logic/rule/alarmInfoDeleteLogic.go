package rulelogic

import (
	"context"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/things/service/udsvr/internal/repo/relationDB"
	"gorm.io/gorm"

	"gitee.com/i-Things/things/service/udsvr/internal/svc"
	"gitee.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmInfoDeleteLogic {
	return &AlarmInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AlarmInfoDeleteLogic) AlarmInfoDelete(in *ud.WithID) (*ud.Empty, error) {
	err := stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		err := relationDB.NewAlarmInfoRepo(tx).Delete(l.ctx, in.Id)
		if err != nil {
			return err
		}
		return relationDB.NewAlarmSceneRepo(tx).DeleteByFilter(l.ctx, relationDB.AlarmSceneFilter{AlarmID: in.Id})
	})

	return &ud.Empty{}, err
}
