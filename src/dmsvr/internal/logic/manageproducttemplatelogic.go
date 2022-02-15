package logic

import (
	"context"
	"database/sql"
	"github.com/go-things/things/shared/errors"
	mysql "github.com/go-things/things/src/dmsvr/internal/repo/mysql"
	"github.com/spf13/cast"
	"time"

	"github.com/go-things/things/src/dmsvr/dm"
	"github.com/go-things/things/src/dmsvr/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ManageProductTemplateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewManageProductTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ManageProductTemplateLogic {
	return &ManageProductTemplateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func UpdateProductTemplate(old *mysql.ProductTemplate, data *dm.ProductTemplate) (isModify bool) {
	defer func() {
		if isModify {
			old.UpdatedTime = sql.NullTime{Valid: true, Time: time.Now()}
		}
	}()
	if data.Template != nil {
		old.Template = data.Template.GetValue()
		isModify = true
	}
	return
}

func (l *ManageProductTemplateLogic) ModifyProductTemplate(in *dm.ManageProductTemplateReq, pt *mysql.ProductTemplate) (*dm.ProductTemplate, error) {
	UpdateProductTemplate(pt, in.Info)
	err := l.svcCtx.ProductTemplate.Update(pt)
	if err != nil {
		l.Errorf("ModifyProductTemplate|ProductTemplate|Update|err=%+v", err)
		return nil, errors.System.AddDetail(err.Error())
	}
	return ToProductTemplate(pt), nil
}

func (l *ManageProductTemplateLogic) AddProductTemplate(in *dm.ManageProductTemplateReq) (*dm.ProductTemplate, error) {
	pi, err := l.svcCtx.ProductInfo.FindOne(in.Info.ProductID)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.Parameter.AddDetail("not find ProductID id:" + cast.ToString(in.Info.ProductID))
		}
		return nil, errors.Database.AddDetail(err.Error())
	}
	pt := &mysql.ProductTemplate{
		ProductID:   pi.ProductID,
		Template:    in.Info.Template.GetValue(),
		CreatedTime: time.Now(),
	}
	l.svcCtx.ProductTemplate.Insert(pt)

	return ToProductTemplate(pt), nil
}

func (l *ManageProductTemplateLogic) InsertProductTemplate(in *dm.ManageProductTemplateReq) (*dm.ProductTemplate, error) {
	pi, err := l.svcCtx.ProductInfo.FindOne(in.Info.ProductID)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.Parameter.AddDetail("not find ProductID id:" + cast.ToString(in.Info.ProductID))
		}
		return nil, errors.Database.AddDetail(err.Error())
	}
	pt := &mysql.ProductTemplate{
		ProductID:   pi.ProductID,
		Template:    in.Info.Template.GetValue(),
		CreatedTime: time.Now(),
	}
	_, err = l.svcCtx.ProductTemplate.Insert(pt)
	if err != nil {
		return nil, errors.Database.AddDetail(err.Error())
	}
	return ToProductTemplate(pt), nil
}

// 产品模板管理
func (l *ManageProductTemplateLogic) ManageProductTemplate(in *dm.ManageProductTemplateReq) (*dm.ProductTemplate, error) {
	l.Infof("ManageProductTemplate|req=%+v", in)
	pt, err := l.svcCtx.ProductTemplate.FindOne(in.Info.ProductID)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.Parameter.AddDetail("not find ProductID id:" + cast.ToString(in.Info.ProductID))
		}
		return nil, errors.System.AddDetail(err.Error())
	}
	return ToProductTemplate(pt), nil
}
