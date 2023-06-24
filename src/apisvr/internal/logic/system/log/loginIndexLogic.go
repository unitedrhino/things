package log

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"github.com/jinzhu/copier"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginIndexLogic {
	return &LoginIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginIndexLogic) LoginIndex(req *types.SysLogLoginIndexReq) (resp *types.SysLogLoginIndexResp, err error) {
	l.Infof("%s req=%v", utils.FuncName(), req)
	var page sys.PageInfo
	copier.Copy(&page, req.Page)
	info, err := l.svcCtx.LogRpc.LoginLogIndex(l.ctx, &sys.LoginLogIndexReq{
		Page:          &page,
		IpAddr:        req.IpAddr,
		LoginLocation: req.LoginLocation,
		Date:          &sys.DateRange{Start: req.DateRange.Start, End: req.DateRange.End},
	})
	if err != nil {
		return nil, err
	}

	var total int64
	total = info.Total

	var logLoginInfo []*types.SysLogLoginIndexData
	logLoginInfo = make([]*types.SysLogLoginIndexData, 0, len(logLoginInfo))

	for _, i := range info.Info {
		logLoginInfo = append(logLoginInfo, &types.SysLogLoginIndexData{
			UserID:        i.UserID,
			UserName:      i.UserName,
			IpAddr:        i.IpAddr,
			LoginLocation: i.LoginLocation,
			Browser:       i.Browser,
			Os:            i.Os,
			Code:          i.Code,
			Msg:           i.Msg,
			CreatedTime:   i.CreatedTime,
		})
	}

	return &types.SysLogLoginIndexResp{List: logLoginInfo, Total: total}, nil
}
