package userlogic

import (
	"github.com/i-Things/things/src/syssvr/internal/repo/mysql"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func UserInfoToPb(ui *mysql.SysUserInfo) *sys.UserInfo {
	return &sys.UserInfo{
		Uid:         ui.Uid,
		UserName:    ui.UserName.String,
		Password:    ui.Password,
		Email:       ui.Email.String,
		Phone:       ui.Phone.String,
		Wechat:      ui.Wechat.String,
		LastIP:      ui.LastIP,
		RegIP:       ui.RegIP,
		NickName:    ui.NickName,
		Role:        ui.Role,
		Sex:         ui.Sex,
		IsAllData:   ui.IsAllData,
		City:        ui.City,
		Country:     ui.Country,
		Province:    ui.Province,
		Language:    ui.Language,
		HeadImgUrl:  ui.HeadImgUrl,
		CreatedTime: ui.CreatedTime.Unix(),
	}
}
