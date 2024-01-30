package usermanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/oss/common"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"github.com/zeromicro/go-zero/core/logx"
)

func UserInfoToPb(ctx context.Context, ui *relationDB.SysUserInfo, svcCtx *svc.ServiceContext) *sys.UserInfo {
	if ui.HeadImg != "" {
		var err error
		ui.HeadImg, err = svcCtx.OssClient.PrivateBucket().SignedGetUrl(ctx, ui.HeadImg, 24*60*60, common.OptionKv{})
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

func ToUserAreaApplyInfos(in []*relationDB.SysUserAreaApply) (ret []*sys.UserAreaApplyInfo) {
	for _, v := range in {
		ret = append(ret, &sys.UserAreaApplyInfo{
			Id:          v.ID,
			UserID:      v.UserID,
			AreaID:      int64(v.AreaID),
			AuthType:    v.AuthType,
			CreatedTime: v.CreatedTime.Unix(),
		})
	}
	return
}
