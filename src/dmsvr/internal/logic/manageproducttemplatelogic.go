package logic

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceTemplate"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/spf13/cast"
	"time"

	"github.com/i-Things/things/src/dmsvr/dm"
	"github.com/i-Things/things/src/dmsvr/internal/svc"

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

func (l *ManageProductTemplateLogic) ModifyProductTemplate(in *dm.ManageProductTemplateReq, pt *mysql.ProductTemplate) (*dm.ProductTemplate, error) {
	//todo 这里需要添加模板的校验
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
		Template:    in.Info.Template,
		CreatedTime: time.Now(),
	}
	tempMode, err := deviceTemplate.NewTemplate([]byte(in.Info.Template))
	if err != nil {
		return nil, errors.Parameter.WithMsg("模板格式不正确").AddDetail(err)
	}
	fmt.Println(tempMode)
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
		Template:    in.Info.Template,
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
			l.AddProductTemplate(in)
		}
		return nil, errors.System.AddDetail(err.Error())
	}
	return l.ModifyProductTemplate(in, pt)
}
