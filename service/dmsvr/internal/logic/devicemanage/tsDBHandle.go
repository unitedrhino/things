package devicemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/share/devices"
)

func UpdateDevice(ctx context.Context, svcCtx *svc.ServiceContext, devs []*devices.Core, affiliation devices.Affiliation) error {
	svcCtx.AbnormalRepo.UpdateDevice(ctx, devs, affiliation)
	return nil
}
