package firmware

import (
	"context"

	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/oss/common"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.OtaFirmwareDelReq) error {
	firmwareReq := &dm.FirmwareInfoDeleteReq{
		FirmwareID: req.FirmwareID,
	}
	deleteResp, err := l.svcCtx.FirmwareM.FirmwareInfoDelete(l.ctx, firmwareReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.FirmwareInfoDelete|req=%v|err=%+v", utils.FuncName(), req, er)
		return er
	}
	//删除附件
	for _, v := range deleteResp.Path {
		l.svcCtx.OssClient.Delete(l.ctx, v, common.OptionKv{})
	}
	return nil
}
