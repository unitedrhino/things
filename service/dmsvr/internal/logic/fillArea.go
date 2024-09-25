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
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func FillAreaDeviceCount(ctx context.Context, svcCtx *svc.ServiceContext, areaIDPaths ...string) error {
	logx.WithContext(ctx).Infof("FillAreaDeviceCount areaIDPaths:%v", areaIDPaths)
	defer utils.Recover(ctx)
	ctx = ctxs.WithRoot(ctx)
	log := logx.WithContext(ctx)
	var idMap = map[int64]struct{}{}
	for _, areaIDPath := range areaIDPaths {
		if areaIDPath == "" || areaIDPath == def.NotClassifiedPath {
			continue
		}
		ids := utils.GetIDPath(areaIDPath)
		var idPath string
		for _, id := range ids {
			idPath += cast.ToString(id) + "-"
			if _, ok := idMap[id]; ok {
				continue
			}
			idMap[id] = struct{}{}
			count, err := relationDB.NewDeviceInfoRepo(ctx).CountByFilter(ctx, relationDB.DeviceFilter{AreaIDPath: idPath})
			if err != nil {
				log.Error(err)
				continue
			}
			_, err = svcCtx.AreaM.AreaInfoUpdate(ctx, &sys.AreaInfo{AreaID: id, DeviceCount: &wrapperspb.Int64Value{Value: count}})
			if err != nil {
				log.Error(err)
				continue
			}
		}
	}

	return nil
}

func FillProjectDeviceCount(ctx context.Context, svcCtx *svc.ServiceContext, projectIDs ...int64) error {
	logx.WithContext(ctx).Infof("FillProjectDeviceCount projectIDs:%v", projectIDs)
	defer utils.Recover(ctx)
	ctx = ctxs.WithRoot(ctx)
	log := logx.WithContext(ctx)
	for _, id := range projectIDs {
		if id <= def.NotClassified {
			continue
		}
		count, err := relationDB.NewDeviceInfoRepo(ctx).CountByFilter(ctx, relationDB.DeviceFilter{ProjectIDs: []int64{id}})
		if err != nil {
			log.Error(err)
			continue
		}
		_, err = svcCtx.ProjectM.ProjectInfoUpdate(ctx, &sys.ProjectInfo{ProjectID: id, DeviceCount: &wrapperspb.Int64Value{Value: count}})
		if err != nil {
			log.Error(err)
			continue
		}
	}

	return nil
}
