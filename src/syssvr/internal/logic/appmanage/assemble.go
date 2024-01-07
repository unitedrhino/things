package appmanagelogic

import (
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func ToAppInfoPo(in *sys.AppInfo) *relationDB.SysAppInfo {
	if in == nil {
		return nil
	}
	return &relationDB.SysAppInfo{
		ID:      in.Id,
		Code:    in.Code,
		Name:    in.Name,
		Type:    in.Type,
		Desc:    utils.ToEmptyString(in.Desc),
		BaseUrl: in.BaseUrl,
		LogoUrl: in.LogoUrl,
	}
}

func ToAppInfoPb(in *relationDB.SysAppInfo) *sys.AppInfo {
	if in == nil {
		return nil
	}
	return &sys.AppInfo{
		Id:      in.ID,
		Code:    in.Code,
		Type:    in.Type,
		Name:    in.Name,
		Desc:    utils.ToRpcNullString(in.Desc),
		BaseUrl: in.BaseUrl,
		LogoUrl: in.LogoUrl,
	}
}

func ToAppInfosPb(in []*relationDB.SysAppInfo) []*sys.AppInfo {
	var ret []*sys.AppInfo
	for _, v := range in {
		ret = append(ret, ToAppInfoPb(v))
	}
	return ret
}
