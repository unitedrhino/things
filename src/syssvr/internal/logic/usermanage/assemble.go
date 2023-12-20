package usermanagelogic

import (
	"github.com/i-Things/things/shared/domain/userDataAuth"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func UserInfoToPb(ui *relationDB.SysUserInfo) *sys.UserInfo {
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
		HeadImgUrl:    ui.HeadImgUrl,
		CreatedTime:   ui.CreatedTime.Unix(),
	}
}

func transAreaPoToPb(po *relationDB.SysUserAuthArea) *sys.UserArea {
	return &sys.UserArea{
		AreaID: int64(po.AreaID),
	}
}

func transProjectPoToPb(po *relationDB.SysUserAuthProject) *sys.UserProject {
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

func DBToAuthAreaDo(area *relationDB.SysUserAuthArea) *userDataAuth.Area {
	if area == nil {
		return nil
	}
	return &userDataAuth.Area{AreaID: int64(area.AreaID)}
}
func DBToAuthAreaDos(areas []*relationDB.SysUserAuthArea) (ret []*userDataAuth.Area) {
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

func DBToAuthProjectDo(area *relationDB.SysUserAuthProject) *userDataAuth.Project {
	if area == nil {
		return nil
	}
	return &userDataAuth.Project{ProjectID: int64(area.ProjectID)}
}
func DBToAuthProjectDos(areas []*relationDB.SysUserAuthProject) (ret []*userDataAuth.Project) {
	if len(areas) == 0 {
		return
	}
	for _, v := range areas {
		ret = append(ret, DBToAuthProjectDo(v))
	}
	return
}
