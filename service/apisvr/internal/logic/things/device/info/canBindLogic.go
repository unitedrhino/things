package info

import (
	"context"
	"strconv"
	"strings"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CanBindLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// canBindProjectIDDetailPrefix 解析 dmsvr 在 DeviceBound 详情中携带的非共享设备项目。
const canBindProjectIDDetailPrefix = "projectID="

// canBindIsShareDetail 解析 dmsvr 在 DeviceBound 详情中携带的共享关系标记。
const canBindIsShareDetail = "isShare=true"

func NewCanBindLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CanBindLogic {
	return &CanBindLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CanBind 检查设备是否可绑定，并提取当前用户已有访问关系时的前端线索。
func (l *CanBindLogic) CanBind(req *types.DeviceInfoCanBindReq) (*types.DeviceInfoCanBindResp, error) {
	_, err := l.svcCtx.DeviceM.DeviceInfoCanBind(l.ctx, utils.Copy[dm.DeviceInfoCanBindReq](req))
	if err != nil {
		if resp, ok := parseCanBindProjectInfo(err); ok {
			return resp, err
		}
		return nil, err
	}
	return nil, nil
}

// parseCanBindProjectInfo 从 DeviceBound 错误详情中解析非共享设备项目和共享关系。
func parseCanBindProjectInfo(err error) (*types.DeviceInfoCanBindResp, bool) {
	if err == nil || !errors.Cmp(err, errors.DeviceBound) {
		return nil, false
	}
	resp := &types.DeviceInfoCanBindResp{}
	hasData := false
	for _, detail := range errors.Fmt(err).Details {
		if detail == canBindIsShareDetail {
			resp.IsShare = true
			hasData = true
			continue
		}
		projectIDText, ok := strings.CutPrefix(detail, canBindProjectIDDetailPrefix)
		if !ok {
			continue
		}
		projectID, parseErr := strconv.ParseInt(projectIDText, 10, 64)
		if parseErr != nil || projectID <= 0 {
			return nil, false
		}
		resp.ProjectID = projectID
		hasData = true
	}
	return resp, hasData
}
