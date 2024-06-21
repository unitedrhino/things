package logic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
)

func FillAreaDeviceCount(ctx context.Context, svcCtx *svc.ServiceContext, areaIDPaths ...string) error {
	ctx = ctxs.WithRoot(ctx)
	log := logx.WithContext(ctx)
	for _, areaIDPath := range areaIDPaths {
		if areaIDPath == "" || areaIDPath == def.NotClassifiedPath {
			continue
		}
		ids := utils.GetIDPath(areaIDPath)
		var idPath string
		for _, id := range ids {
			idPath += cast.ToString(id) + "-"
			count, err := relationDB.NewDeviceInfoRepo(ctx).CountByFilter(ctx, relationDB.DeviceFilter{AreaIDPath: idPath})
			if err != nil {
				log.Error(err)
				continue
			}
			_, err = svcCtx.AreaM.AreaInfoUpdate(ctx, &sys.AreaInfo{AreaID: id, DeviceCount: count})
			if err != nil {
				log.Error(err)
				continue
			}
		}
	}

	return nil
}
