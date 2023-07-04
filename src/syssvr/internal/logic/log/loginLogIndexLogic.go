package loglogic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/repo/mysql"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogIndexLogic {
	return &LoginLogIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogIndexLogic) LoginLogIndex(in *sys.LoginLogIndexReq) (*sys.LoginLogIndexResp, error) {
	resp, total, err := l.svcCtx.LogModel.LoginLogIndex(l.ctx, &mysql.LoginLogFilter{
		Page:          &def.PageInfo{Page: in.Page.Page, Size: in.Page.Size},
		IpAddr:        in.IpAddr,
		LoginLocation: in.LoginLocation,
		Data: &mysql.DateRange{
			Start: in.Date.Start,
			End:   in.Date.End,
		},
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}

	info := make([]*sys.LoginLogIndexData, 0, len(resp))
	for _, v := range resp {
		info = append(info, &sys.LoginLogIndexData{
			UserID:        v.UserID,
			UserName:      v.UserName,
			IpAddr:        v.IpAddr,
			LoginLocation: v.LoginLocation,
			Browser:       v.Browser,
			Os:            v.Os,
			Code:          v.Code,
			Msg:           v.Msg,
			CreatedTime:   v.CreatedTime.Unix(),
		})
	}

	return &sys.LoginLogIndexResp{Info: info, Total: total}, nil
}
