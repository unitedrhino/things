package logic

import (
	"context"
	"time"

	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func FillAreaGroupCount(ctx context.Context, svcCtx *svc.ServiceContext, areaID int64) error {
	logx.WithContext(ctx).Infof("FillAreaDeviceCount areaID:%v", areaID)
	defer utils.Recover(ctx)
	ctx = ctxs.WithRoot(ctx)
	log := logx.WithContext(ctx)
	if areaID <= def.NotClassified {
		return nil
	}
	count, err := stores.WithNoDebug(ctx, relationDB.NewGroupInfoRepo).CountByFilter(ctx, relationDB.GroupInfoFilter{AreaID: areaID})
	if err != nil {
		log.Error(err)
		return err
	}
	area, err := svcCtx.AreaCache.GetData(ctx, areaID)
	if err != nil {
		log.Error(err)
	}
	if area.GroupCount.GetValue() == count {
		return nil
	}
	_, err = svcCtx.AreaM.AreaInfoUpdate(ctx, &sys.AreaInfo{ProjectID: area.ProjectID, AreaID: areaID, GroupCount: &wrapperspb.Int64Value{Value: count}})
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func DirectFillAreaDeviceCount(ctx context.Context, svcCtx *svc.ServiceContext, delay time.Duration, areas ...*sys.AreaInfo) error {
	logx.WithContext(ctx).Infof("FillAreaDeviceCount delay:%v len:%v", delay, len(areas))
	defer utils.Recover(ctx)
	var startTime = time.Now()
	ctx = ctxs.WithRoot(ctx)
	log := logx.WithContext(ctx)
	var idMap = map[int64]struct{}{}
	var updateCount int
	for _, area := range areas {
		if area.AreaID <= def.NotClassified || area.AreaIDPath == "" || area.AreaIDPath == def.NotClassifiedPath {
			continue
		}
		ids := utils.GetIDPath(area.AreaIDPath)
		var idPath string
		for _, id := range ids {
			idPath += cast.ToString(id) + "-"
			if _, ok := idMap[id]; ok {
				continue
			}
			if delay != 0 {
				time.Sleep(delay)
			}
			idMap[id] = struct{}{}
			count, err := stores.WithNoDebug(ctx, relationDB.NewDeviceInfoRepo).CountByFilter(ctx, relationDB.DeviceFilter{AreaIDPath: idPath})
			if err != nil {
				log.Error(err)
				continue
			}
			if area.DeviceCount.GetValue() == count {
				continue
			}
			updateCount++
			_, err = svcCtx.AreaM.AreaInfoUpdate(ctx, &sys.AreaInfo{ProjectID: area.ProjectID, AreaID: id, DeviceCount: &wrapperspb.Int64Value{Value: count}})
			if err != nil {
				log.Error(err)
				continue
			}
		}
	}
	log.Infof("FillAreaDeviceCount change:%v use:%v", updateCount, time.Since(startTime))

	return nil
}

func FillAreaDeviceCount(ctx context.Context, svcCtx *svc.ServiceContext, areas ...*sys.AreaInfo) error {
	areaChan <- areas
	return nil
}

var projectIDChan = make(chan []int64, 100000)
var areaChan = make(chan []*sys.AreaInfo, 100000)

func Init(svcCtx *svc.ServiceContext) {
	utils.Go(context.Background(), func() {
		tick := time.Tick(time.Second)
		execProjectIDs := make([]int64, 0, 500)
		execAreas := make([]*sys.AreaInfo, 0, 1000)
		for {
			select {
			case _ = <-tick:
				if len(execProjectIDs) > 0 {
					batchSize := len(execProjectIDs)
					if batchSize > 500 { //控制每秒执行的速率
						batchSize = 500
					}
					newProjectIDs := execProjectIDs[:batchSize]
					execProjectIDs = execProjectIDs[batchSize:] //清空切片
					utils.Go(context.Background(), func() {
						ctx := ctxs.WithRoot(context.Background())
						DirectFillProjectDeviceCount(ctx, svcCtx, 0, newProjectIDs...)
					})
				}
				if len(execAreas) > 0 {
					batchSize := len(execAreas)
					if batchSize > 500 { //控制每秒执行的速率
						batchSize = 500
					}
					newAreas := execAreas[:batchSize]
					execAreas = execAreas[batchSize:] //清空切片
					utils.Go(context.Background(), func() {
						ctx := ctxs.WithRoot(context.Background())
						DirectFillAreaDeviceCount(ctx, svcCtx, 0, newAreas...)
					})
				}
			case p := <-projectIDChan:
				execProjectIDs = append(execProjectIDs, p...)
			case a := <-areaChan:
				execAreas = append(execAreas, a...)
			}
		}
	})
}
func DirectFillProjectDeviceCount(ctx context.Context, svcCtx *svc.ServiceContext, delay time.Duration, projectIDs ...int64) error {
	logx.WithContext(ctx).Infof("FillProjectDeviceCount delay:%v len:%v", delay, len(projectIDs))
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
		if delay != 0 {
			time.Sleep(delay)
		}
		idMap[id] = struct{}{}
		count, err := stores.WithNoDebug(ctx, relationDB.NewDeviceInfoRepo).CountByFilter(ctx, relationDB.DeviceFilter{ProjectIDs: []int64{id}})
		if err != nil {
			log.Error(err)
			continue
		}
		pi, err := svcCtx.ProjectCache.GetData(ctx, id)
		if err != nil {
			log.Error(err)
			continue
		}
		if pi.DeviceCount.GetValue() == count {
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
