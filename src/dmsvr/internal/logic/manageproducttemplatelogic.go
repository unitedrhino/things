package logic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
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
	l.Infof("ManageProductTemplate|ModifyProductTemplate|ProductID:%v", in.Info.ProductID)
	newTempMode, err := deviceTemplate.ValidateWithFmt([]byte(in.Info.Template))
	if err != nil {
		return nil, err
	}
	newT, err := deviceTemplate.NewTemplate(newTempMode)
	if err != nil {
		return nil, err
	}
	oldT, err := deviceTemplate.NewTemplate([]byte(pt.Template))
	if err != nil {
		l.Errorf("%s new old template failure,err:%v,old:%v", utils.FuncName(), err, pt.Template)
		return nil, err
	}
	err = deviceTemplate.CheckModify(oldT, newT)
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.DeviceDataRepo.ModifyProduct(l.ctx, oldT, newT, in.Info.ProductID); err != nil {
		l.Errorf("%s ModifyProduct failure,err:%v", utils.FuncName(), err)
		return nil, errors.Database.AddDetail(err)
	}
	pt.Template = string(newTempMode)
	err = l.svcCtx.ProductTemplate.Update(pt)
	if err != nil {
		l.Errorf("ModifyProductTemplate|ProductTemplate|Update|err=%+v", err)
		return nil, errors.System.AddDetail(err.Error())
	}
	return ToProductTemplate(pt), nil
}

func (l *ManageProductTemplateLogic) AddProductTemplate(in *dm.ManageProductTemplateReq) (*dm.ProductTemplate, error) {
	l.Infof("ManageProductTemplate|AddProductTemplate|ProductID:%v", in.Info.ProductID)
	pi, err := l.svcCtx.ProductInfo.FindOne(in.Info.ProductID)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.Parameter.AddDetail("not find ProductID id:" + cast.ToString(in.Info.ProductID))
		}
		return nil, errors.Database.AddDetail(err.Error())
	}
	newTempMode, err := deviceTemplate.ValidateWithFmt([]byte(in.Info.Template))
	if err != nil {
		return nil, err
	}
	t, err := deviceTemplate.NewTemplate(newTempMode)
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.DeviceDataRepo.InitProduct(l.ctx, t, in.Info.ProductID); err != nil {
		l.Errorf("%s InitProduct failure,err:%v", utils.FuncName(), err)
		return nil, errors.Database.AddDetail(err)
	}

	pt := &mysql.ProductTemplate{
		ProductID:   pi.ProductID,
		Template:    string(newTempMode),
		CreatedTime: time.Now(),
	}
	l.svcCtx.ProductTemplate.Insert(pt)
	return ToProductTemplate(pt), nil
}

// 产品模板管理
func (l *ManageProductTemplateLogic) ManageProductTemplate(in *dm.ManageProductTemplateReq) (*dm.ProductTemplate, error) {
	l.Infof("ManageProductTemplate|req=%+v", in)
	pt, err := l.svcCtx.ProductTemplate.FindOne(in.Info.ProductID)
	if err != nil {
		if err == mysql.ErrNotFound {
			return l.AddProductTemplate(in)
		}
		return nil, errors.System.AddDetail(err.Error())
	}
	return l.ModifyProductTemplate(in, pt)
}
