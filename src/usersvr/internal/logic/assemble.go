package logic

import (
	"github.com/i-Things/things/src/usersvr/internal/repo/mysql"
	"github.com/i-Things/things/src/usersvr/pb/user"
)

func UserCoreToDb() {

}
func UserCoreToPb(core *mysql.UserCore) *user.UserCore {
	return &user.UserCore{
		Uid:         core.Uid,
		UserName:    core.UserName,
		Password:    core.Password,
		Email:       core.Email,
		Phone:       core.Phone,
		Wechat:      core.Wechat,
		LastIP:      core.LastIP,
		RegIP:       core.RegIP,
		CreatedTime: core.CreatedTime.Unix(),
		Status:      core.Status,
		Role:        core.AuthorityId,
	}
}

func UserInfoToPb(ui *mysql.UserInfo) *user.UserInfo {
	return &user.UserInfo{
		Uid:        ui.Uid,
		UserName:   ui.UserName,
		NickName:   ui.NickName,
		InviterUid: ui.InviterUid,
		InviterId:  ui.InviterId,
		Sex:        ui.Sex,
		City:       ui.City,
		Country:    ui.Country,
		Province:   ui.Province,
		Language:   ui.Language,
		HeadImgUrl: ui.HeadImgUrl,
		CreateTime: ui.CreatedTime.Unix(),
	}
}

func UserInfoToDb(ui *user.UserInfo) *mysql.UserInfo {
	return &mysql.UserInfo{
		Uid:        ui.Uid,
		UserName:   ui.UserName,
		NickName:   ui.NickName,
		InviterUid: ui.InviterUid,
		InviterId:  ui.InviterId,
		Sex:        ui.Sex,
		City:       ui.City,
		Country:    ui.Country,
		Province:   ui.Province,
		Language:   ui.Language,
		HeadImgUrl: ui.HeadImgUrl,
		//CreatedTime: sql.NullTime{Valid: true, Time: time.Now()},
	}
}
