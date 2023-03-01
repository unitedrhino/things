package loglogic

import (
	"context"
	"database/sql"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/domain/log"
	"github.com/i-Things/things/src/syssvr/internal/repo/mysql"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperLogCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOperLogCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperLogCreateLogic {
	return &OperLogCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OperLogCreateLogic) OperLogCreate(in *sys.OperLogCreateReq) (*sys.Response, error) {
	//OperUserName 用uid查用户表获得
	resUser, err := l.svcCtx.UserInfoModel.FindOne(l.ctx, in.Uid)
	if err != nil {
		return nil, errors.Database.AddMsgf("UserInfoModel.FindOne is err, uid:%ld", in.Uid)
	}
	//OperName，BusinessType 用Route查接口管理表获得
	resApi, err := l.svcCtx.ApiModel.FindOneByRoute(l.ctx, in.Route)
	if err != nil {
		return nil, errors.Database.AddMsgf("ApiModel.FindOneByRoute is err, url:%s", in.Route)
	}
	if resApi.BusinessType != log.OptQuery {
		_, err := l.svcCtx.LogOperModel.Insert(l.ctx, &mysql.SysOperLog{
			OperUid:      in.Uid,
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
			return nil, errors.Database.AddMsgf("LogOperModel.Insert is err:%s", err.Error())
		}
	}

	return &sys.Response{}, nil
}
