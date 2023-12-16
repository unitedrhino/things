package loglogic

import (
	"context"
	"database/sql"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/domain/log"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperLogCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	OlDB *relationDB.OperLogRepo
	UiDB *relationDB.UserInfoRepo
	AiDB *relationDB.ApiInfoRepo
}

func NewOperLogCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperLogCreateLogic {
	return &OperLogCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		OlDB:   relationDB.NewOperLogRepo(ctx),
		UiDB:   relationDB.NewUserInfoRepo(ctx),
		AiDB:   relationDB.NewApiInfoRepo(ctx),
	}
}

func (l *OperLogCreateLogic) OperLogCreate(in *sys.OperLogCreateReq) (*sys.Response, error) {
	//OperUserName 用uid查用户表获得
	resUser, err := l.UiDB.FindOne(l.ctx, in.UserID)
	if err != nil {
		return nil, errors.Database.AddMsgf("UserInfoModel.FindOne is err, UserID:%ld", in.UserID)
	}
	//OperName，BusinessType 用Route查接口管理表获得
	resApi, err := l.AiDB.FindOneByFilter(l.ctx, relationDB.ApiInfoFilter{Route: in.Route})
	if err != nil {
		return nil, errors.Database.AddMsgf("ApiModel.FindOneByRoute is err, url:%s", in.Route)
	}
	if resApi.BusinessType != log.OptQuery {
		err := l.OlDB.Insert(l.ctx, &relationDB.SysOperLog{
			AppCode:      in.AppCode,
			OperUserID:   in.UserID,
			OperUserName: resUser.UserName.String,
			OperName:     resApi.Name,
			BusinessType: resApi.BusinessType,
			Uri:          in.Uri,
			OperIpAddr:   in.OperIpAddr,
			OperLocation: in.OperLocation,
			Req:          sql.NullString{String: in.Req, Valid: true},
			Resp:         sql.NullString{String: in.Resp, Valid: true},
			Code:         in.Code,
			Msg:          in.Msg,
		})
		if err != nil {
			return nil, err
		}
	}

	return &sys.Response{}, nil
}
