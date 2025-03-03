package logic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"time"
)

func FillAreaGroupCount(ctx context.Context, svcCtx *svc.ServiceContext, areaID int64) error {
	logx.WithContext(ctx).Infof("FillAreaDeviceCount areaID:%v", areaID)
	defer utils.Recover(ctx)
	ctx = ctxs.WithRoot(ctx)
	log := logx.WithContext(ctx)
	if areaID <= def.NotClassified {
		return nil
	}
	count, err := relationDB.NewGroupInfoRepo(ctx).CountByFilter(ctx, relationDB.GroupInfoFilter{AreaID: areaID})
	if err != nil {
		log.Error(err)
		return err
	}
	area, err := svcCtx.AreaCache.GetData(ctx, areaID)
	if err != nil {
		log.Error(err)
	}
	_, err = svcCtx.AreaM.AreaInfoUpdate(ctx, &sys.AreaInfo{ProjectID: area.ProjectID, AreaID: areaID, GroupCount: &wrapperspb.Int64Value{Value: count}})
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func fillAreaDeviceCount(ctx context.Context, svcCtx *svc.ServiceContext, areaIDPaths ...string) error {
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
			area, err := svcCtx.AreaCache.GetData(ctx, id)
			if err != nil {
				log.Error(err)
				continue
			}
			_, err = svcCtx.AreaM.AreaInfoUpdate(ctx, &sys.AreaInfo{ProjectID: area.ProjectID, AreaID: id, DeviceCount: &wrapperspb.Int64Value{Value: count}})
			if err != nil {
				log.Error(err)
				continue
			}
		}
	}

	return nil
}

func FillAreaDeviceCount(ctx context.Context, svcCtx *svc.ServiceContext, areaIDPaths ...string) error {
	areaIDPathChan <- areaIDPaths
	return nil
}

var projectIDChan chan []int64
var areaIDPathChan chan []string

func Init(svcCtx *svc.ServiceContext) {
	projectIDChan = make(chan []int64, 100000)
	areaIDPathChan = make(chan []string, 100000)
	utils.Go(context.Background(), func() {
		tick := time.Tick(time.Second)
		execProjectIDs := make([]int64, 0, 500)
		execAreaIDPaths := make([]string, 0, 1000)
		for {
			select {
			case _ = <-tick:
				if len(execProjectIDs) > 0 {
					batchSize := len(execProjectIDs)
					if batchSize > 100 { //控制每秒执行的速率
						batchSize = 100
					}
					newProjectIDs := execProjectIDs[:batchSize]
					execProjectIDs = execProjectIDs[batchSize:] //清空切片
					utils.Go(context.Background(), func() {
						ctx := ctxs.WithRoot(context.Background())
						fillProjectDeviceCount(ctx, svcCtx, newProjectIDs...)
					})
				}
				if len(execAreaIDPaths) > 0 {
					batchSize := len(execProjectIDs)
					if batchSize > 100 { //控制每秒执行的速率
						batchSize = 100
					}
					newAreaIDPaths := execAreaIDPaths[:batchSize]
					execProjectIDs = execProjectIDs[batchSize:] //清空切片
					execAreaIDPaths = execAreaIDPaths[0:0]      //清空切片
					utils.Go(context.Background(), func() {
						ctx := ctxs.WithRoot(context.Background())
						fillAreaDeviceCount(ctx, svcCtx, newAreaIDPaths...)
					})
				}
			case p := <-projectIDChan:
				execProjectIDs = append(execProjectIDs, p...)
			case a := <-areaIDPathChan:
				execAreaIDPaths = append(execAreaIDPaths, a...)
			}
		}
	})
}
func fillProjectDeviceCount(ctx context.Context, svcCtx *svc.ServiceContext, projectIDs ...int64) error {
	logx.WithContext(ctx).Infof("FillProjectDeviceCount projectIDs:%v", projectIDs)
	defer utils.Recover(ctx)
	ctx = ctxs.WithRoot(ctx)
	log := logx.WithContext(ctx)
	var idMap = map[int64]struct{}{}
	for _, id := range projectIDs {
		if id <= def.NotClassified {
			continue
		}
		if _, ok := idMap[id]; ok {
			continue
		}
		idMap[id] = struct{}{}
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

func FillProjectDeviceCount(ctx context.Context, svcCtx *svc.ServiceContext, projectIDs ...int64) error {
	projectIDChan <- projectIDs
	return nil
}
