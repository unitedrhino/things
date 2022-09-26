package info

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/jinzhu/copier"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func toTagsMap(tags []*types.Tag) map[string]string {
	if tags == nil {
		return nil
	}
	tagMap := make(map[string]string, len(tags))
	for _, tag := range tags {
		tagMap[tag.Key] = tag.Value
	}
	return tagMap
}

func toTagsType(tags map[string]string) (retTag []*types.Tag) {
	for k, v := range tags {
		retTag = append(retTag, &types.Tag{
			Key:   k,
			Value: v,
		})
	}
	return
}

func (l *IndexLogic) Index(req *types.GroupInfoIndexReq) (resp *types.GroupInfoIndexResp, err error) {
	var page dm.PageInfo
	copier.Copy(&page, req.Page)

	res, err := l.svcCtx.DeviceG.GroupInfoIndex(l.ctx, &dm.GroupInfoIndexReq{
		Page:      &page,
		GroupName: req.GroupName,
		Tags:      toTagsMap(req.Tags),
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.GroupInfoIndex req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	glist := make([]*types.GroupInfo, 0, len(res.List))
	for _, v := range res.List {
		glist = append(glist, &types.GroupInfo{
			GroupName:   v.GroupName,
			GroupID:     v.GroupID,
			ParentID:    v.ParentID,
			CreatedTime: v.CreatedTime,
			Desc:        v.Desc,
			Tags:        toTagsType(v.Tags),
		})
	}

	glistAll := make([]*types.GroupInfo, 0, len(res.ListAll))
	for _, v := range res.ListAll {
		glistAll = append(glistAll, &types.GroupInfo{
			GroupName:   v.GroupName,
			GroupID:     v.GroupID,
			ParentID:    v.ParentID,
			CreatedTime: v.CreatedTime,
			Desc:        v.Desc,
			Tags:        toTagsType(v.Tags),
		})
	}

	return &types.GroupInfoIndexResp{
		List:    glist,
		Total:   res.Total,
		ListAll: glistAll,
	}, nil
}
