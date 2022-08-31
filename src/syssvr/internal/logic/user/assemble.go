package userlogic

import (
	"github.com/i-Things/things/src/syssvr/internal/repo/mysql"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"time"
)

func UserCoreToDb() {

}

func UserInfoToPb(ui *mysql.UserInfo) *sys.UserInfo {
	return &sys.UserInfo{
		Uid:         ui.Uid,
		UserName:    ui.UserName,
		Password:    ui.Password,
		Email:       ui.Email,
		Phone:       ui.Phone,
		Wechat:      ui.Wechat,
		LastIP:      ui.LastIP,
		RegIP:       ui.RegIP,
		Role:        ui.Role,
		NickName:    ui.NickName,
		Sex:         ui.Sex,
		City:        ui.City,
		Country:     ui.Country,
		Province:    ui.Province,
		Language:    ui.Language,
		HeadImgUrl:  ui.HeadImgUrl,
		CreatedTime: ui.CreatedTime.Unix(),
	}
}

func UserInfoToDb(ui *sys.UserInfo) *mysql.UserInfo {
	return &mysql.UserInfo{
		Uid:         ui.Uid,
		UserName:    ui.UserName,
		Password:    ui.Password,
		Email:       ui.Email,
		Phone:       ui.Phone,
		Wechat:      ui.Wechat,
		LastIP:      ui.LastIP,
		RegIP:       ui.RegIP,
		Role:        ui.Role,
		NickName:    ui.NickName,
		Sex:         ui.Sex,
		City:        ui.City,
		Country:     ui.Country,
		Province:    ui.Province,
		Language:    ui.Language,
		HeadImgUrl:  ui.HeadImgUrl,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
}
