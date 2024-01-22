package otafirmwaremanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/shared/utils/cast"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"gorm.io/gorm"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaFirmwareDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB  *relationDB.ProductInfoRepo
	OfDB  *relationDB.OtaFirmwareRepo
	OffDB *relationDB.OtaFirmwareFileRepo
}

func NewOtaFirmwareDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaFirmwareDeleteLogic {
	return &OtaFirmwareDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProductInfoRepo(ctx),
		OfDB:   relationDB.NewOtaFirmwareRepo(ctx),
		OffDB:  relationDB.NewOtaFirmwareFileRepo(ctx),
	}
}

// 删除升级包
func (l *OtaFirmwareDeleteLogic) OtaFirmwareDelete(in *dm.OtaFirmwareDeleteReq) (*dm.Response, error) {
	_, err := l.OfDB.FindOneByFilter(l.ctx, relationDB.OtaFirmwareFilter{FirmwareID: in.FirmwareId})
	if errors.Cmp(err, errors.NotFind) {
		l.Errorf("not find firmware id:" + cast.ToString(in.FirmwareId))
		return nil, err
	} else if err != nil {
		return nil, err
	}
	//开启事务
	db := stores.GetCommonConn(l.ctx)
	err = db.Transaction(func(tx *gorm.DB) error {
		//删除升级包文件
		err = l.OffDB.Delete(l.ctx, in.FirmwareId)
		if err != nil {
			l.Errorf("%s.DeleteOTAFirmwareFile err=%v", utils.FuncName(), err)
			tx.Rollback()
			return errors.System.AddDetail(err)
		}
		//删除升级包
		err = l.OfDB.Delete(l.ctx, in.FirmwareId)
		if err != nil {
			l.Errorf("%s.DeleteOTAFirmware err=%v", utils.FuncName(), err)
			tx.Rollback()
			return errors.System.AddDetail(err)
		}
		// 如果所有操作成功，提交事务
		return nil
	})
	if err != nil {
		l.Errorf("failed to commit transaction: %v", err)
		return nil, errors.System.AddDetail(err)
	}
	return &dm.Response{}, nil
}
