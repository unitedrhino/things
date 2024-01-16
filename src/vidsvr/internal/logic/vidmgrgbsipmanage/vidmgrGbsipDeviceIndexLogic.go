package vidmgrgbsipmanagelogic

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/src/vidsvr/internal/common"
	"github.com/i-Things/things/src/vidsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/pb/vid"

	"github.com/zeromicro/go-zero/core/logx"
)

type VidmgrGbsipDeviceIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVidmgrGbsipDeviceIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VidmgrGbsipDeviceIndexLogic {
	return &VidmgrGbsipDeviceIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取GB28181设备列表
func (l *VidmgrGbsipDeviceIndexLogic) VidmgrGbsipDeviceIndex(in *vid.VidmgrGbsipDeviceIndexReq) (*vid.VidmgrGbsipDeviceIndexResp, error) {
	// todo: add your logic here and delete this line
	deviceRepo := relationDB.NewVidmgrDevicesRepo(l.ctx)
	filter := relationDB.VidmgrDevicesFilter{
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
	info := make([]*vid.VidmgrGbsipDevice, 0, len(di))
	for _, v := range di {
		info = append(info, common.ToVidmgrGbsipDeviceRpc(v))
	}
	fmt.Printf("----airgens-----VidmgrGbsipDeviceIndex:")
	return &vid.VidmgrGbsipDeviceIndexResp{List: info, Total: size}, nil
}
