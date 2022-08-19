package userlogic

import (
	"github.com/i-Things/things/src/syssvr/internal/repo/mysql"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func UserCoreToDb() {

}
func UserCoreToPb(core *mysql.UserCore) *sys.UserCore {
	return &sys.UserCore{
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

func UserInfoToPb(ui *mysql.UserInfo) *sys.UserInfo {
	return &sys.UserInfo{
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

func UserInfoToDb(ui *sys.UserInfo) *mysql.UserInfo {
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
