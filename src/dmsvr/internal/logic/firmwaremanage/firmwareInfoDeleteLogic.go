package firmwaremanagelogic

import (
	"context"

	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type FirmwareInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFirmwareInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FirmwareInfoDeleteLogic {
	return &FirmwareInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FirmwareInfoDeleteLogic) FirmwareInfoDelete(in *dm.FirmwareInfoDeleteReq) (*dm.FirmwareInfoDeleteResp, error) {
	var fDB = relationDB.NewOtaFirmwareRepo(l.ctx)
	var fileDB = relationDB.NewOtaFirmwareFileRepo(l.ctx)
	//删除DB数据
	err := fDB.Delete(l.ctx, in.FirmwareID)
	if err != nil {
		l.Errorf("DelFirmware|FirmwareInfo|Delete|err=%+v", err)
		return nil, errors.System.AddDetail(err)
	}
	fi, err := fileDB.FindByFilter(l.ctx, relationDB.OtaFirmwareFileFilter{
		FirmwareID: in.FirmwareID,
	}, &def.PageInfo{Size: 20, Page: 1})
	if err != nil {
		return nil, err
	}
	//删除升级任务
	relationDB.NewOtaTaskRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.OtaTaskFilter{FirmwareID: in.FirmwareID})
	relationDB.NewOtaTaskDevicesRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.OtaTaskDevicesFilter{FirmwareID: in.FirmwareID})
	//TODO 删除附件，或者 修改
	var m = make([]string, len(fi))
	for k, v := range fi {
		m[k] = v.FilePath
	}
	return &dm.FirmwareInfoDeleteResp{Path: m}, nil
}
