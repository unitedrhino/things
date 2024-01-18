package usermanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"gorm.io/gorm"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAreaApplyDealLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserAreaApplyDealLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAreaApplyDealLogic {
	return &UserAreaApplyDealLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserAreaApplyDealLogic) UserAreaApplyDeal(in *sys.UserAreaApplyDealReq) (*sys.Response, error) {
	if !in.IsApprove {
		err := relationDB.NewUserAreaApplyRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.UserAreaApplyFilter{IDs: in.Ids})
		return &sys.Response{}, err
	}
	db := stores.GetTenantConn(l.ctx)
	err := db.Transaction(func(tx *gorm.DB) error {
		uaa := relationDB.NewUserAreaApplyRepo(tx)
		ua := relationDB.NewUserAreaRepo(tx)
		uaas, err := uaa.FindByFilter(l.ctx, relationDB.UserAreaApplyFilter{IDs: in.Ids}, nil)
		if err != nil {
			return err
		}
		var uas []*relationDB.SysUserArea
		for _, v := range uaas {
			uas = append(uas, &relationDB.SysUserArea{
				UserID:    v.UserID,
				ProjectID: int64(v.ProjectID),
				AreaID:    int64(v.AreaID),
				AuthType:  v.AuthType,
			})
		}
		err = ua.MultiInsert(l.ctx, uas)
		if err != nil {
			return err
		}
		err = uaa.DeleteByFilter(l.ctx, relationDB.UserAreaApplyFilter{IDs: in.Ids})
		if err != nil {
			return err
		}
		return nil
	})

	return &sys.Response{}, err
}
