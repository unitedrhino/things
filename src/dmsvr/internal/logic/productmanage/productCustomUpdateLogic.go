package productmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductCustomUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductCustomUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductCustomUpdateLogic {
	return &ProductCustomUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProductCustomUpdateLogic) ProductCustomUpdate(in *dm.ProductCustom) (*dm.Response, error) {
	pi, err := l.svcCtx.ProductCustom.FindOneByProductID(l.ctx, in.ProductID)
	if err != nil {
		if err == mysql.ErrNotFound {
			if in.ScriptLang == 0 {
				in.ScriptLang = 1
			}
			_, err = l.svcCtx.ProductCustom.Insert(l.ctx, &mysql.DmProductCustom{
				ProductID:       in.ProductID,
				ScriptLang:      in.ScriptLang,
				CustomTopic:     utils.AnyToNullString(in.CustomTopic),
				TransformScript: utils.AnyToNullString(in.TransformScript),
			})
			if err != nil {
				return nil, errors.Database.AddDetail(err)
			}
			return &dm.Response{}, nil
		}
		return nil, errors.Database.AddDetail(err)
	}
	if in.TransformScript != nil {
		pi.TransformScript = utils.AnyToNullString(in.TransformScript)
	}
	if in.ScriptLang != 0 {
		pi.ScriptLang = in.ScriptLang
	}
	if in.CustomTopic != nil {
		pi.CustomTopic = utils.AnyToNullString(in.CustomTopic)
	}
	err = l.svcCtx.ProductCustom.Update(l.ctx, pi)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	err = l.svcCtx.DataUpdate.ProductCustomUpdate(l.ctx, &events.DeviceUpdateInfo{ProductID: in.ProductID})
	if err != nil {
		return nil, err
	}
	return &dm.Response{}, nil
}
