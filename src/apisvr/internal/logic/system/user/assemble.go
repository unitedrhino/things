package user

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func GetNullVal(val *wrappers.StringValue) *string {
	if val == nil {
		return nil
	}
	return &val.Value
}

func UserCoreToApi(core *sys.UserCore) *types.UserCore {
	return &types.UserCore{
		Uid:         core.Uid,
		UserName:    core.UserName,
		Email:       core.Email,
		Phone:       core.Phone,
		Wechat:      core.Wechat,
		LastIP:      core.LastIP,
		RegIP:       core.RegIP,
		CreatedTime: core.CreatedTime,
		Status:      core.Status,
		Role:        core.Role,
	}
}

func UserInfoToApi(ui *sys.UserInfo) *types.UserInfo {
	return &types.UserInfo{
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
		CreateTime: ui.CreateTime,
	}
}

//func UserInfoFullToApi(ui *user.UserInfo) *types.UserIndexResp {
//	return &types.UserIndexResp{
//		Uid:        ui.Uid,
//		UserName:   ui.UserName,
//		NickName:   ui.NickName,
//		InviterUid: ui.InviterUid,
//		InviterId:  ui.InviterId,
//		Sex:        ui.Sex,
//		City:       ui.City,
//		Country:    ui.Country,
//		Province:   ui.Province,
//		Language:   ui.Language,
//		HeadImgUrl: ui.HeadImgUrl,
//		CreateTime: ui.CreateTime,
//	}
//}
