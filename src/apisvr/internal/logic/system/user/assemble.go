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

func UserInfoToApi(ui *sys.UserInfo) *types.UserInfo {
	return &types.UserInfo{
		UserID:      ui.UserID,
		UserName:    ui.UserName,
		Email:       ui.Email,
		Phone:       ui.Phone,
		Wechat:      ui.Wechat,
		LastIP:      ui.LastIP,
		RegIP:       ui.RegIP,
		Role:        ui.Role,
		NickName:    ui.NickName,
		Sex:         ui.Sex,
		IsAllData:   ui.IsAllData,
		City:        ui.City,
		Country:     ui.Country,
		Province:    ui.Province,
		Language:    ui.Language,
		HeadImgUrl:  ui.HeadImgUrl,
		CreatedTime: ui.CreatedTime,
	}
}
func UserInfoToRpc(ui *types.UserInfo) *sys.UserInfo {
	return &sys.UserInfo{
		UserID:      ui.UserID,
		UserName:    ui.UserName,
		Email:       ui.Email,
		Phone:       ui.Phone,
		Wechat:      ui.Wechat,
		LastIP:      ui.LastIP,
		RegIP:       ui.RegIP,
		Role:        ui.Role,
		NickName:    ui.NickName,
		Sex:         ui.Sex,
		IsAllData:   ui.IsAllData,
		City:        ui.City,
		Country:     ui.Country,
		Province:    ui.Province,
		Language:    ui.Language,
		HeadImgUrl:  ui.HeadImgUrl,
		CreatedTime: ui.CreatedTime,
	}
}
