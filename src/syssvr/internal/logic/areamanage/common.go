package areamanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"strings"
)

func checkProject(ctx context.Context, tx *gorm.DB, projectID int64) (*relationDB.SysProjectInfo, error) {
	if projectID == 0 {
		projectID = ctxs.GetUserCtx(ctx).ProjectID
	}
	po, err := relationDB.NewProjectInfoRepo(tx).FindOne(ctx, projectID)
	if err == nil {
		return po, nil
	}
	if errors.Cmp(err, errors.NotFind) {
		return nil, nil
	}
	return nil, err
}

func checkArea(ctx context.Context, tx *gorm.DB, areaID int64) (*relationDB.SysAreaInfo, error) {
	po, err := relationDB.NewAreaInfoRepo(tx).FindOne(ctx, areaID, nil)
	if err == nil {
		return po, nil
	}
	if errors.Cmp(err, errors.NotFind) {
		return nil, nil
	}
	return nil, err
}

func checkParentArea(ctx context.Context, tx *gorm.DB, parentAreaID int64) (*relationDB.SysAreaInfo, error) {
	//检查父级区域是否存在
	parentAreaPo, err := checkArea(ctx, tx, parentAreaID)
	if err != nil {
		return nil, errors.Fmt(err).WithMsg("检查区域出错")
	} else if parentAreaPo == nil {
		return nil, errors.Parameter.AddDetail(parentAreaID).WithMsg("检查区域不存在")
	}
	return parentAreaPo, nil
}

func addSubAreaIDs(ctx context.Context, tx *gorm.DB, po *relationDB.SysAreaInfo, subAreaID int64) error {
	po.ChildrenAreaIDs = append(po.ChildrenAreaIDs, subAreaID)
	ids := GetIDPath(po.AreaIDPath)
	if len(ids) > 1 {
		ids = ids[:len(ids)-1]
		areas, err := relationDB.NewAreaInfoRepo(tx).FindByFilter(ctx, relationDB.AreaInfoFilter{AreaIDs: ids}, nil)
		if err != nil {
			return err
		}
		for _, v := range areas {
			v.ChildrenAreaIDs = append(v.ChildrenAreaIDs, subAreaID)
			err := relationDB.NewAreaInfoRepo(tx).Update(ctx, v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func subSubAreaIDs(ctx context.Context, tx *gorm.DB, po *relationDB.SysAreaInfo, subAreaID int64) error {
	po.ChildrenAreaIDs = utils.SliceDelete(po.ChildrenAreaIDs, subAreaID)
	ids := GetIDPath(po.AreaIDPath)
	if len(ids) > 1 {
		ids = ids[:len(ids)-1]
		areas, err := relationDB.NewAreaInfoRepo(tx).FindByFilter(ctx, relationDB.AreaInfoFilter{AreaIDs: ids}, nil)
		if err != nil {
			return err
		}
		for _, v := range areas {
			v.ChildrenAreaIDs = utils.SliceDelete(v.ChildrenAreaIDs, subAreaID)
			err := relationDB.NewAreaInfoRepo(tx).Update(ctx, v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GetIDPath(areaIDPath string) (ret []int64) {
	ids := strings.Split(areaIDPath, "-")
	for _, v := range ids {
		if v != "" {
			ret = append(ret, cast.ToInt64(v))
		}
	}
	return ret
}
func GetNamePath(areaNamePath string) (ret []string) {
	ids := strings.Split(areaNamePath, "-")
	for _, v := range ids {
		if v != "" {
			ret = append(ret, v)
		}
	}
	return ret
}
