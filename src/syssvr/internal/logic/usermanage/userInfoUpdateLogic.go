package usermanagelogic

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/oss"
	"github.com/i-Things/things/shared/oss/common"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UiDB *relationDB.UserInfoRepo
}

func NewUserInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoUpdateLogic {
	return &UserInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UiDB:   relationDB.NewUserInfoRepo(ctx),
	}
}

func (l *UserInfoUpdateLogic) UserInfoUpdate(in *sys.UserInfo) (*sys.Response, error) {
	ui, err := l.UiDB.FindOneByFilter(l.ctx, relationDB.UserInfoFilter{UserIDs: []int64{in.UserID}, WithRoles: true})
	if err != nil {
		l.Errorf("%s.FindOne UserID=%d err=%v", utils.FuncName(), in.UserID, err)
		return nil, err
	}
	if in.UserName != "" {
		ui.UserName = sql.NullString{String: in.UserName, Valid: true}
	}
	if in.NickName != "" {
		ui.NickName = in.NickName
	}

	//性別有效才賦值，否則使用旧值
	if ui.Sex == def.Unknown {
		ui.Sex = def.Male
	} else {
		ui.Sex = in.Sex
	}

	//设置数据超管
	if in.IsAllData == 1 || in.IsAllData == 2 {
		ui.IsAllData = in.IsAllData
	}
	if in.Role != 0 && in.Role != ui.Role {
		ui.Role = in.Role
	}
	if in.IsUpdateHeadImg == def.True && in.HeadImg != "" {
		if ui.HeadImg != "" {
			err := l.svcCtx.OssClient.PrivateBucket().Delete(l.ctx, ui.HeadImg, common.OptionKv{})
			if err != nil {
				l.Errorf("Delete file err path:%v,err:%v", ui.HeadImg, err)
			}
		}
		nwePath := oss.GenFilePath(l.ctx, l.svcCtx.Config.Name, oss.BusinessUserManage, oss.SceneUserInfo, fmt.Sprintf("%d/%s", ui.UserID, oss.GetFileNameWithPath(in.HeadImg)))
		path, err := l.svcCtx.OssClient.PrivateBucket().CopyFromTempBucket(in.HeadImg, nwePath)
		if err != nil {
			return nil, errors.System.AddDetail(err)
		}
		ui.HeadImg = path
	}

	if in.Role != 0 { //默认角色只能修改为授权的角色
		for _, r := range ui.Roles {
			if r.RoleID == in.Role {
				ui.Role = in.Role
			}
		}
	}

	err = l.UiDB.Update(l.ctx, ui)
	if err != nil {
		l.Errorf("%s.Update ui=%v err=%v", utils.FuncName(), ui, err)
		return nil, err
	}
	l.Infof("%s.modified usersvr info = %+v", utils.FuncName(), ui)

	return &sys.Response{}, nil
}
