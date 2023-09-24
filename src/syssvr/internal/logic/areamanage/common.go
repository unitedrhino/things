package areamanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
)

func checkProject(ctx context.Context, productID int64) (*relationDB.SysProjectInfo, error) {
	po, err := relationDB.NewProjectInfoRepo(ctx).FindOne(ctx, productID)
	if err == nil {
		return po, nil
	}
	if errors.Cmp(err, errors.NotFind) {
		return nil, nil
	}
	return nil, err
}

func checkArea(ctx context.Context, areaID int64) (*relationDB.SysAreaInfo, error) {
	po, err := relationDB.NewAreaInfoRepo(ctx).FindOne(ctx, areaID, nil)
	if err == nil {
		return po, nil
	}
	if errors.Cmp(err, errors.NotFind) {
		return nil, nil
	}
	return nil, err
}

func checkParentArea(ctx context.Context, parentAreaID int64, checkDevice bool) (*relationDB.SysAreaInfo, error) {
	//检查父级区域是否存在
	parentAreaPo, err := checkArea(ctx, parentAreaID)
	if err != nil {
		return nil, errors.Database.AddDetail(err).WithMsg("检查区域出错")
	} else if parentAreaPo == nil {
		return nil, errors.Parameter.AddDetail(parentAreaID).WithMsg("检查区域不存在")
	}
	return parentAreaPo, nil
}
