package usermanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/domain/userDataAuth"
	"github.com/i-Things/things/shared/oss/common"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"github.com/zeromicro/go-zero/core/logx"
)

func UserInfoToPb(ctx context.Context, ui *relationDB.SysTenantUserInfo, svcCtx *svc.ServiceContext) *sys.UserInfo {
	if ui.HeadImg != "" {
		var err error
		ui.HeadImg, err = svcCtx.OssClient.PrivateBucket().SignedGetUrl(ctx, ui.HeadImg, 24*60, common.OptionKv{})
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.SignedGetUrl err:%v", utils.FuncName(), err)
		}
	}
	return &sys.UserInfo{
		UserID:        ui.UserID,
		UserName:      ui.UserName.String,
		Email:         ui.Email.String,
		Phone:         ui.Phone.String,
		WechatUnionID: ui.WechatUnionID.String,
		LastIP:        ui.LastIP,
		RegIP:         ui.RegIP,
		NickName:      ui.NickName,
		Role:          ui.Role,
		Sex:           ui.Sex,
		IsAllData:     ui.IsAllData,
		City:          ui.City,
		Country:       ui.Country,
		Province:      ui.Province,
		Language:      ui.Language,
		HeadImg:       ui.HeadImg,
		CreatedTime:   ui.CreatedTime.Unix(),
	}
}

func transAreaPoToPb(po *relationDB.SysUserArea) *sys.UserArea {
	return &sys.UserArea{
		AreaID: int64(po.AreaID),
	}
}

func transProjectPoToPb(po *relationDB.SysUserProject) *sys.UserProject {
	return &sys.UserProject{
		ProjectID: int64(po.ProjectID),
	}
}

func ToAuthAreaDo(area *sys.UserArea) *userDataAuth.Area {
	if area == nil {
		return nil
	}
	return &userDataAuth.Area{AreaID: area.AreaID}
}
func ToAuthAreaDos(areas []*sys.UserArea) (ret []*userDataAuth.Area) {
	if len(areas) == 0 {
		return
	}
	for _, v := range areas {
		ret = append(ret, ToAuthAreaDo(v))
	}
	return
}

func DBToAuthAreaDo(area *relationDB.SysUserArea) *userDataAuth.Area {
	if area == nil {
		return nil
	}
	return &userDataAuth.Area{AreaID: int64(area.AreaID)}
}
func DBToAuthAreaDos(areas []*relationDB.SysUserArea) (ret []*userDataAuth.Area) {
	if len(areas) == 0 {
		return
	}
	for _, v := range areas {
		ret = append(ret, DBToAuthAreaDo(v))
	}
	return
}

func ToAuthProjectDo(area *sys.UserProject) *userDataAuth.Project {
	if area == nil {
		return nil
	}
	return &userDataAuth.Project{ProjectID: area.ProjectID}
}
func ToAuthProjectDos(areas []*sys.UserProject) (ret []*userDataAuth.Project) {
	if len(areas) == 0 {
		return
	}
	for _, v := range areas {
		ret = append(ret, ToAuthProjectDo(v))
	}
	return
}

func DBToAuthProjectDo(area *relationDB.SysUserProject) *userDataAuth.Project {
	if area == nil {
		return nil
	}
	return &userDataAuth.Project{ProjectID: int64(area.ProjectID)}
}
func DBToAuthProjectDos(areas []*relationDB.SysUserProject) (ret []*userDataAuth.Project) {
	if len(areas) == 0 {
		return
	}
	for _, v := range areas {
		ret = append(ret, DBToAuthProjectDo(v))
	}
	return
}

func ToUserAreaApplyInfos(in []*relationDB.SysUserAreaApply) (ret []*sys.UserAreaApplyInfo) {
	for _, v := range in {
		ret = append(ret, &sys.UserAreaApplyInfo{
			Id:          v.ID,
			AreaID:      int64(v.AreaID),
			AuthType:    v.AuthType,
			CreatedTime: v.CreatedTime.Unix(),
		})
	}
	return
}
