package logic

import (
	"database/sql"
	"github.com/i-Things/things/shared/utils/cast"
	"github.com/i-Things/things/src/usersvr/model"
	"github.com/i-Things/things/src/usersvr/user"
	"time"
)

func UserCoreToDb() {

}
func UserCoreToPb(core *model.UserCore) *user.UserCore {
	return &user.UserCore{
		Uid:         core.Uid,
		UserName:    core.UserName,
		Password:    core.Password,
		Email:       core.Email,
		Phone:       core.Phone,
		Wechat:      core.Wechat,
		LastIP:      core.LastIP,
		RegIP:       core.RegIP,
		CreatedTime: cast.ToInt64(core.CreatedTime),
		Status:      core.Status,
	}
}

func UserInfoToPb(ui *model.UserInfo) *user.UserInfo {
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
		CreateTime: cast.ToInt64(ui.CreatedTime),
	}
}

func UserInfoToDb(ui *user.UserInfo) *model.UserInfo {
	return &model.UserInfo{
		Uid:         ui.Uid,
		UserName:    ui.UserName,
		NickName:    ui.NickName,
		InviterUid:  ui.InviterUid,
		InviterId:   ui.InviterId,
		Sex:         ui.Sex,
		City:        ui.City,
		Country:     ui.Country,
		Province:    ui.Province,
		Language:    ui.Language,
		HeadImgUrl:  ui.HeadImgUrl,
		CreatedTime: sql.NullTime{Valid: true, Time: time.Now()},
	}
}
