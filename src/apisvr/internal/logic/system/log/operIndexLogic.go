package log

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOperIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperIndexLogic {
	return &OperIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OperIndexLogic) OperIndex(req *types.SysLogOperIndexReq) (resp *types.SysLogOperIndexResp, err error) {
	l.Infof("%s req=%v", utils.FuncName(), req)
	info, err := l.svcCtx.LogRpc.OperLogIndex(l.ctx, &sys.OperLogIndexReq{
		Page:         logic.ToSysPageRpc(req.Page),
		AppCode:      req.AppCode,
		OperName:     req.OperName,
		OperUserName: req.OperUserName,
		BusinessType: req.BusinessType,
	})
	if err != nil {
		return nil, err
	}

	var total int64
	total = info.Total

	var logOperInfo []*types.SysLogOperInfo
	logOperInfo = make([]*types.SysLogOperInfo, 0, len(logOperInfo))

	for _, i := range info.List {
		logOperInfo = append(logOperInfo, &types.SysLogOperInfo{
			AppCode:      i.AppCode,
			UserID:       i.UserID,
			OperUserName: i.OperUserName,
			OperName:     i.OperName,
			BusinessType: i.BusinessType,
			Uri:          i.Uri,
			OperIpAddr:   i.OperIpAddr,
			OperLocation: i.OperLocation,
			Req:          i.Req,
			Resp:         i.Resp,
			Code:         i.Code,
			Msg:          i.Msg,
			CreatedTime:  i.CreatedTime,
		})
	}

	return &types.SysLogOperIndexResp{List: logOperInfo, Total: total}, nil
}
