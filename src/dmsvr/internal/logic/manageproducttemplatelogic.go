package logic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/dm"
	"github.com/i-Things/things/src/dmsvr/internal/domain/thing"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/spf13/cast"

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

func (l *ManageProductTemplateLogic) ModifyProductTemplate(in *dm.ManageProductTemplateReq, oldT *thing.Template) (*dm.ProductTemplate, error) {
	l.Infof("ManageProductTemplate|ModifyProductTemplate|ProductID:%v", in.Info.ProductID)
	newT, err := thing.ValidateWithFmt([]byte(in.Info.Template))
	if err != nil {
		return nil, err
	}
	err = thing.CheckModify(oldT, newT)
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.DeviceDataRepo.ModifyProduct(l.ctx, oldT, newT, in.Info.ProductID); err != nil {
		l.Errorf("%s ModifyProduct failure,err:%v", utils.FuncName(), err)
		return nil, errors.Database.AddDetail(err)
	}
	err = l.svcCtx.TemplateRepo.Update(l.ctx, in.Info.ProductID, newT)
	if err != nil {
		l.Errorf("ModifyProductTemplate|ProductTemplate|Update|err=%+v", err)
		return nil, errors.System.AddDetail(err.Error())
	}
	err = l.svcCtx.DataUpdate.TempModelUpdate(l.ctx, &thing.TemplateInfo{ProductID: in.Info.ProductID})
	if err != nil {
		return nil, err
	}
	pt, err := l.svcCtx.TemplateRepo.GetTemplateInfo(l.ctx, in.Info.ProductID)
	if err != nil {
		return nil, err
	}
	return ToProductTemplate(pt), nil
}

func (l *ManageProductTemplateLogic) AddProductTemplate(in *dm.ManageProductTemplateReq) (*dm.ProductTemplate, error) {
	l.Infof("ManageProductTemplate|AddProductTemplate|ProductID:%v", in.Info.ProductID)
	_, err := l.svcCtx.ProductInfo.FindOne(l.ctx, in.Info.ProductID)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.Parameter.AddDetail("not find ProductID id:" + cast.ToString(in.Info.ProductID))
		}
		return nil, errors.Database.AddDetail(err.Error())
	}
	t, err := thing.ValidateWithFmt([]byte(in.Info.Template))
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.HubLogRepo.InitProduct(
		l.ctx, in.Info.ProductID); err != nil {
		l.Errorf("%s|DeviceLogRepo|InitProduct| failure,err:%v", utils.FuncName(), err)
		return nil, errors.Database.AddDetail(err)
	}
	if err := l.svcCtx.DeviceDataRepo.InitProduct(l.ctx, t, in.Info.ProductID); err != nil {
		l.Errorf("%s|DeviceDataRepo|InitProduct| failure,err:%v", utils.FuncName(), err)
		return nil, errors.Database.AddDetail(err)
	}
	err = l.svcCtx.TemplateRepo.Insert(l.ctx, in.Info.ProductID, t)
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.DataUpdate.TempModelUpdate(l.ctx, &thing.TemplateInfo{ProductID: in.Info.ProductID})
	if err != nil {
		return nil, err
	}
	pt, err := l.svcCtx.TemplateRepo.GetTemplateInfo(l.ctx, in.Info.ProductID)
	if err != nil {
		return nil, err
	}
	return ToProductTemplate(pt), err
}

// 产品模板管理
func (l *ManageProductTemplateLogic) ManageProductTemplate(in *dm.ManageProductTemplateReq) (*dm.ProductTemplate, error) {
	l.Infof("ManageProductTemplate|req=%+v", in)
	pt, err := l.svcCtx.TemplateRepo.GetTemplate(l.ctx, in.Info.ProductID)
	if err != nil {
		if err == mysql.ErrNotFound {
			return l.AddProductTemplate(in)
		}
		return nil, errors.System.AddDetail(err.Error())
	}
	return l.ModifyProductTemplate(in, pt)
}
