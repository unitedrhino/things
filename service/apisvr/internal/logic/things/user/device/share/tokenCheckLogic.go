package share

import (
	"context"
	"strings"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type TokenCheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 公开查询批量分享 Token 状态
func NewTokenCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TokenCheckLogic {
	return &TokenCheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TokenCheckLogic) TokenCheck(req *types.UserDeviceShareTokenCheckReq) (*types.UserDeviceShareTokenCheckResp, error) {
	shareToken := strings.TrimSpace(req.ShareToken)
	if shareToken == "" {
		return nil, errors.Parameter.WithMsg("shareToken不能为空")
	}
	ret, err := l.svcCtx.UserDevice.UserDeivceShareMultiIndex(l.ctx, &dm.UserDeviceShareMultiToken{ShareToken: shareToken})
	if err != nil {
		if isTokenCheckExpiredOrConsumed(err) {
			return InvalidTokenCheckResp(), nil
		}
		return nil, err
	}
	return ToTokenCheckResp(ret), nil
}

func isTokenCheckExpiredOrConsumed(err error) bool {
	if errors.Cmp(err, errors.NotFind) {
		return true
	}
	errText := err.Error()
	return strings.Contains(errText, "分享已过期") || strings.Contains(errText, "不存在")
}
