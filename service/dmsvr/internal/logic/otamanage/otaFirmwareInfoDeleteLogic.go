package otamanagelogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/oss/common"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/spf13/cast"
	"gorm.io/gorm"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaFirmwareInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB  *relationDB.ProductInfoRepo
	OfDB  *relationDB.OtaFirmwareInfoRepo
	OffDB *relationDB.OtaFirmwareFileRepo
}

func NewOtaFirmwareInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaFirmwareInfoDeleteLogic {
	return &OtaFirmwareInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProductInfoRepo(ctx),
		OfDB:   relationDB.NewOtaFirmwareInfoRepo(ctx),
		OffDB:  relationDB.NewOtaFirmwareFileRepo(ctx),
	}
}

// 删除升级包
func (l *OtaFirmwareInfoDeleteLogic) OtaFirmwareInfoDelete(in *dm.WithID) (*dm.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	l.ctx = ctxs.WithRoot(l.ctx)
	_, err := l.OfDB.FindOneByFilter(l.ctx, relationDB.OtaFirmwareInfoFilter{ID: in.Id})
	if errors.Cmp(err, errors.NotFind) {
		l.Errorf("not find firmware id:" + cast.ToString(in.Id))
		return nil, err
	} else if err != nil {
		return nil, err
	}
	var (
		fs []*relationDB.DmOtaFirmwareFile
	)

	//开启事务
	db := stores.GetCommonConn(l.ctx)
	err = db.Transaction(func(tx *gorm.DB) error {
		//删除升级包文件
		offDb := relationDB.NewOtaFirmwareFileRepo(tx)
		fs, err = offDb.FindByFilter(l.ctx, relationDB.OtaFirmwareFileFilter{FirmwareID: in.Id}, nil)
		if err != nil {
			return err
		}
		err = offDb.Delete(l.ctx, in.Id)
		if err != nil {
			l.Errorf("%s.DeleteOTAFirmwareFile err=%v", utils.FuncName(), err)
			return err
		}
		ofDb := relationDB.NewOtaFirmwareInfoRepo(tx)
		//删除升级包
		err = ofDb.Delete(l.ctx, in.Id)
		if err != nil {
			l.Errorf("%s.DeleteOTAFirmware err=%v", utils.FuncName(), err)
			return err
		}
		// 如果所有操作成功，提交事务
		return nil
	})
	if err != nil {
		l.Errorf("failed to commit transaction: %v", err)
		return nil, err
	}
	for _, v := range fs {
		err := l.svcCtx.OssClient.PrivateBucket().Delete(l.ctx, v.FilePath, common.OptionKv{})
		if err != nil {
			l.Error(err)
		}
	}
	return &dm.Empty{}, nil
}
