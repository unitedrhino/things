package apilogic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/repo/mysql"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApiIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiIndexLogic {
	return &ApiIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApiIndexLogic) ApiIndex(in *sys.ApiIndexReq) (*sys.ApiIndexResp, error) {
	resp, total, err := l.svcCtx.ApiInfoModel.Index(l.ctx, &mysql.ApiFilter{
		Page:   &def.PageInfo{Page: in.Page.Page, Size: in.Page.Size},
		Route:  in.Route,
		Method: in.Method,
		Group:  in.Group,
		Name:   in.Name,
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}

	info := make([]*sys.ApiData, 0, len(resp))
	for _, v := range resp {
		info = append(info, &sys.ApiData{
			Id:           v.Id,
			Route:        v.Route,
			Method:       v.Method,
			Group:        v.Group,
			Name:         v.Name,
			BusinessType: v.BusinessType,
		})
	}

	return &sys.ApiIndexResp{
		Total: total,
		List:  info,
	}, nil
}
