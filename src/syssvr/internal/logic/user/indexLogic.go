package userlogic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/repo/mysql"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IndexLogic) Index(in *sys.UserIndexReq) (*sys.UserIndexResp, error) {
	l.Infof("%s req=%+v", utils.FuncName(), in)

	ucs, total, err := l.svcCtx.UserModel.Index(&mysql.UserIndexReq{
		Page:     &def.PageInfo{Page: in.Page.Page, Size: in.Page.Size},
		UserName: in.UserName,
		Phone:    in.Phone,
		Email:    in.Email,
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	info := make([]*sys.UserInfo, 0, len(ucs))
	for _, uc := range ucs {
		info = append(info, UserInfoToPb(uc))
	}
	return &sys.UserIndexResp{
		Info:  info,
		Total: total,
	}, nil

}
