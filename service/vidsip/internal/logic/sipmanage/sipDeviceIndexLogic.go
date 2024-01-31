package sipmanagelogic

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/shared/def"
	"github.com/i-Things/things/service/vidsip/internal/logic/common"
	db "github.com/i-Things/things/service/vidsip/internal/repo/relationDB"
	"github.com/i-Things/things/service/vidsip/internal/svc"
	"github.com/i-Things/things/service/vidsip/pb/sip"

	"github.com/zeromicro/go-zero/core/logx"
)

type SipDeviceIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSipDeviceIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SipDeviceIndexLogic {
	return &SipDeviceIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取GB28181设备列表
func (l *SipDeviceIndexLogic) SipDeviceIndex(in *sip.SipDevIndexReq) (*sip.SipDevIndexResp, error) {
	// todo: add your logic here and delete this line
	deviceRepo := db.NewSipDevicesRepo(l.ctx)
	filter := db.SipDevicesFilter{
		DeviceIDs: in.DeviceIDs,
	}
	fmt.Printf("----airgens-----VidmgrGbsipDeviceIndex:")
	size, err := deviceRepo.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}
	di, err := deviceRepo.FindByFilter(l.ctx, filter, common.ToPageInfoWithDefault(in.Page, &def.PageInfo{
		Page: 1, Size: 20,
		Orders: []def.OrderBy{{"created_time", def.OrderDesc}, {"device_id", def.OrderDesc}},
	}))
	if err != nil {
		return nil, err
	}
	info := make([]*sip.SipDevice, 0, len(di))
	for _, v := range di {
		info = append(info, common.ToSipDeviceRpc(v))
	}
	fmt.Printf("----airgens-----VidmgrGbsipDeviceIndex:")

	return &sip.SipDevIndexResp{List: info, Total: size}, nil
}
