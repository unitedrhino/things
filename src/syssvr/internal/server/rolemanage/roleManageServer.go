// Code generated by goctl. DO NOT EDIT.
// Source: sys.proto

package server

import (
	"context"

	"github.com/i-Things/things/src/syssvr/internal/logic/rolemanage"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

type RoleManageServer struct {
	svcCtx *svc.ServiceContext
	sys.UnimplementedRoleManageServer
}

func NewRoleManageServer(svcCtx *svc.ServiceContext) *RoleManageServer {
	return &RoleManageServer{
		svcCtx: svcCtx,
	}
}

func (s *RoleManageServer) RoleInfoCreate(ctx context.Context, in *sys.RoleInfo) (*sys.Response, error) {
	l := rolemanagelogic.NewRoleInfoCreateLogic(ctx, s.svcCtx)
	return l.RoleInfoCreate(in)
}

func (s *RoleManageServer) RoleInfoIndex(ctx context.Context, in *sys.RoleInfoIndexReq) (*sys.RoleInfoIndexResp, error) {
	l := rolemanagelogic.NewRoleInfoIndexLogic(ctx, s.svcCtx)
	return l.RoleInfoIndex(in)
}

func (s *RoleManageServer) RoleInfoUpdate(ctx context.Context, in *sys.RoleInfo) (*sys.Response, error) {
	l := rolemanagelogic.NewRoleInfoUpdateLogic(ctx, s.svcCtx)
	return l.RoleInfoUpdate(in)
}

func (s *RoleManageServer) RoleInfoDelete(ctx context.Context, in *sys.ReqWithID) (*sys.Response, error) {
	l := rolemanagelogic.NewRoleInfoDeleteLogic(ctx, s.svcCtx)
	return l.RoleInfoDelete(in)
}

func (s *RoleManageServer) RoleMenuIndex(ctx context.Context, in *sys.RoleMenuIndexReq) (*sys.RoleMenuIndexResp, error) {
	l := rolemanagelogic.NewRoleMenuIndexLogic(ctx, s.svcCtx)
	return l.RoleMenuIndex(in)
}

func (s *RoleManageServer) RoleMenuMultiUpdate(ctx context.Context, in *sys.RoleMenuMultiUpdateReq) (*sys.Response, error) {
	l := rolemanagelogic.NewRoleMenuMultiUpdateLogic(ctx, s.svcCtx)
	return l.RoleMenuMultiUpdate(in)
}

func (s *RoleManageServer) RoleAppIndex(ctx context.Context, in *sys.RoleAppIndexReq) (*sys.RoleAppIndexResp, error) {
	l := rolemanagelogic.NewRoleAppIndexLogic(ctx, s.svcCtx)
	return l.RoleAppIndex(in)
}

func (s *RoleManageServer) RoleAppMultiUpdate(ctx context.Context, in *sys.RoleAppMultiUpdateReq) (*sys.Response, error) {
	l := rolemanagelogic.NewRoleAppMultiUpdateLogic(ctx, s.svcCtx)
	return l.RoleAppMultiUpdate(in)
}

func (s *RoleManageServer) RoleApiAuth(ctx context.Context, in *sys.RoleApiAuthReq) (*sys.Response, error) {
	l := rolemanagelogic.NewRoleApiAuthLogic(ctx, s.svcCtx)
	return l.RoleApiAuth(in)
}

func (s *RoleManageServer) RoleApiMultiUpdate(ctx context.Context, in *sys.RoleApiMultiUpdateReq) (*sys.Response, error) {
	l := rolemanagelogic.NewRoleApiMultiUpdateLogic(ctx, s.svcCtx)
	return l.RoleApiMultiUpdate(in)
}

func (s *RoleManageServer) RoleApiIndex(ctx context.Context, in *sys.RoleApiIndexReq) (*sys.RoleApiIndexResp, error) {
	l := rolemanagelogic.NewRoleApiIndexLogic(ctx, s.svcCtx)
	return l.RoleApiIndex(in)
}