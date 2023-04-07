package productmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductScriptUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductScriptUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductScriptUpdateLogic {
	return &ProductScriptUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProductScriptUpdateLogic) ProductScriptUpdate(in *dm.ProductScript) (*dm.Response, error) {
	pi, err := l.svcCtx.ProductScript.FindOneByProductID(l.ctx, in.ProductID)
	if err != nil {
		if err == mysql.ErrNotFound {
			_, err = l.svcCtx.ProductScript.Insert(l.ctx, &mysql.DmProductScript{
				ProductID: in.ProductID,
				Script:    in.Script,
				Lang:      in.Lang,
			})
			if err != nil {
				return nil, errors.Database.AddDetail(err)
			}
			return &dm.Response{}, nil
		}
		return nil, errors.Database.AddDetail(err)
	}
	if pi.Script == in.Script {
		return &dm.Response{}, nil
	}
	pi.Script = in.Script
	pi.Lang = in.Lang
	err = l.svcCtx.ProductScript.Update(l.ctx, pi)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	err = l.svcCtx.DataUpdate.ProductScriptUpdate(l.ctx, &events.DataUpdateInfo{ProductID: in.ProductID})
	if err != nil {
		return nil, err
	}
	return &dm.Response{}, nil
}
