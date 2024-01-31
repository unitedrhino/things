package firmware

import (
	"context"

	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.OtaFirmwareCreateReq) error {
	firmwareReq := dm.Firmware{
		Name:       req.Name,
		ProductID:  req.ProductID,
		Version:    req.Version,
		IsDiff:     req.IsDiff,
		SignMethod: req.SignMethod,
	}
	firmwareReq.SignMethod = "md5" //先固定为md5
	if req.Desc != nil {
		firmwareReq.Desc = &wrappers.StringValue{
			Value: *req.Desc,
		}
	}
	if req.ExtData != nil {
		firmwareReq.ExtData = &wrappers.StringValue{
			Value: *req.ExtData,
		}
	}
	if len(req.Files) > 0 {
		for _, v := range req.Files {
			firmwareReq.Files = append(firmwareReq.Files, &dm.OtaFirmwareFile{
				Name:     v.Name,
				FilePath: v.Filepath,
			})
		}
	}
	FirmwareInfo, err := l.svcCtx.FirmwareM.FirmwareInfoCreate(l.ctx, &firmwareReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.ManageDevice|req=%v|err=%+v", utils.FuncName(), req, er)
		return er
	}
	l.svcCtx.FileChan <- FirmwareInfo.FirmwareID
	return nil
}
