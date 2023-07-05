package devicegrouplogic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDb *relationDB.ProductInfoRepo
}

func NewGroupInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupInfoIndexLogic {
	return &GroupInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDb:   relationDB.NewProductInfoRepo(ctx),
	}
}

// 获取分组信息列表
func (l *GroupInfoIndexLogic) GroupInfoIndex(in *dm.GroupInfoIndexReq) (*dm.GroupInfoIndexResp, error) {
	ros, total, err := l.svcCtx.GroupDB.Index(l.ctx, &mysql.GroupFilter{
		Page:      &def.PageInfo{Page: in.Page.Page, Size: in.Page.Size},
		GroupName: in.GroupName,
		ParentID:  in.ParentID,
		Tags:      in.Tags,
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	info := make([]*dm.GroupInfo, 0, len(ros))
	//filterProductID :=
	productFilter := relationDB.ProductFilter{}
	for _, ro := range ros {
		productFilter.ProductIDs = append(productFilter.ProductIDs, ro.ProductID)
		info = append(info, &dm.GroupInfo{
			GroupID:     ro.GroupID,
			ParentID:    ro.ParentID,
			ProjectID:   ro.ProjectID,
			GroupName:   ro.GroupName,
			ProductID:   ro.ProductID,
			Desc:        ro.Desc,
			CreatedTime: ro.CreatedTime,
			Tags:        in.Tags,
		})
	}
	productList, _ := l.PiDb.FindByFilter(l.ctx, productFilter, nil)
	productMap := make(map[string]string, len(productList))
	for _, p := range productList {
		productMap[p.ProductID] = p.ProductName
	}
	for k, v := range productList {
		productList[k].ProductName = productMap[v.ProductID]
	}
	rosAll, err := l.svcCtx.GroupDB.IndexAll(l.ctx, &mysql.GroupFilter{
		Page:      &def.PageInfo{Page: in.Page.Page, Size: in.Page.Size},
		GroupName: in.GroupName,
		ParentID:  in.ParentID,
		Tags:      in.Tags,
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	infoAll := make([]*dm.GroupInfo, 0, len(rosAll))
	for _, ro := range rosAll {
		infoAll = append(infoAll, &dm.GroupInfo{
			GroupID:     ro.GroupID,
			ParentID:    ro.ParentID,
			ProjectID:   ro.ProjectID,
			GroupName:   ro.GroupName,
			ProductName: productMap[ro.ProductID],
			ProductID:   ro.ProductID,
			Desc:        ro.Desc,
			CreatedTime: ro.CreatedTime,
			Tags:        in.Tags,
		})
	}

	return &dm.GroupInfoIndexResp{List: info, Total: total, ListAll: infoAll}, nil
}
