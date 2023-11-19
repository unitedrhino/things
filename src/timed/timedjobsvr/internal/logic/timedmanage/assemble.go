package timedmanagelogic

import (
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/src/timed/internal/repo/relationDB"
	"github.com/i-Things/things/src/timed/timedjobsvr/pb/timedjob"
)

func ToTaskGroupPo(in *timedjob.TaskGroup) *relationDB.TimedTaskGroup {
	if in == nil {
		return nil
	}
	return &relationDB.TimedTaskGroup{
		Code:     in.Code,
		Name:     in.Name,
		Type:     in.Type,
		SubType:  in.SubType,
		Priority: in.Priority,
		Env:      in.Env,
		Config:   in.Config,
	}
}
func ToTaskGroupPb(in *relationDB.TimedTaskGroup) *timedjob.TaskGroup {
	if in == nil {
		return nil
	}
	return &timedjob.TaskGroup{
		Code:     in.Code,
		Name:     in.Name,
		Type:     in.Type,
		SubType:  in.SubType,
		Priority: in.Priority,
		Env:      in.Env,
		Config:   in.Config,
	}
}

func ToTaskGroupPbs(in []*relationDB.TimedTaskGroup) (ret []*timedjob.TaskGroup) {
	for _, v := range in {
		ret = append(ret, ToTaskGroupPb(v))
	}
	return
}

func ToTaskInfoPbs(in []*relationDB.TimedTaskInfo) (ret []*timedjob.TaskInfo) {
	for _, v := range in {
		ret = append(ret, ToTaskInfoPb(v))
	}
	return
}

func ToTaskInfoPb(in *relationDB.TimedTaskInfo) *timedjob.TaskInfo {
	if in == nil {
		return nil
	}
	return &timedjob.TaskInfo{
		GroupCode: in.GroupCode,
		Type:      in.Type,
		Name:      in.Name,
		Code:      in.Code,
		Params:    in.Params,
		CronExpr:  in.CronExpr,
		Status:    in.Status,
		Priority:  in.Priority,
	}
}

func ToTaskInfoPo(in *timedjob.TaskInfo) *relationDB.TimedTaskInfo {
	if in == nil {
		return nil
	}
	return &relationDB.TimedTaskInfo{
		GroupCode: in.GroupCode,
		Type:      in.Type,
		Name:      in.Name,
		Code:      in.Code,
		Params:    in.Params,
		CronExpr:  in.CronExpr,
		Status:    in.Status,
		Priority:  in.Priority,
	}
}

func ToPageInfo(info *timedjob.PageInfo, defaultOrders ...def.OrderBy) *def.PageInfo {
	if info == nil {
		return nil
	}

	var orders = defaultOrders
	if infoOrders := info.GetOrders(); len(infoOrders) > 0 {
		orders = make([]def.OrderBy, 0, len(infoOrders))
		for _, infoOd := range infoOrders {
			if infoOd.GetFiled() != "" {
				orders = append(orders, def.OrderBy{infoOd.GetFiled(), infoOd.GetSort()})
			}
		}
	}

	return &def.PageInfo{
		Page:   info.GetPage(),
		Size:   info.GetSize(),
		Orders: orders,
	}
}

func ToPageInfoWithDefault(info *timedjob.PageInfo, defau *def.PageInfo) *def.PageInfo {
	if page := ToPageInfo(info); page == nil {
		return defau
	} else {
		if page.Page == 0 {
			page.Page = defau.Page
		}
		if page.Size == 0 {
			page.Size = defau.Size
		}
		if len(page.Orders) == 0 {
			page.Orders = defau.Orders
		}
		return page
	}
}
