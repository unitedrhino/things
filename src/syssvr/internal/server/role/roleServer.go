// Code generated by goctl. DO NOT EDIT!
// Source: sys.proto

package server

import (
	"context"

	"github.com/i-Things/things/src/syssvr/internal/logic/role"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

type RoleServer struct {
	svcCtx *svc.ServiceContext
	sys.UnimplementedRoleServer
}

func NewRoleServer(svcCtx *svc.ServiceContext) *RoleServer {
	return &RoleServer{
		svcCtx: svcCtx,
	}
}

func (s *RoleServer) RoleCreate(ctx context.Context, in *sys.RoleCreateReq) (*sys.Response, error) {
	l := rolelogic.NewRoleCreateLogic(ctx, s.svcCtx)
	return l.RoleCreate(in)
}

func (s *RoleServer) RoleIndex(ctx context.Context, in *sys.RoleIndexReq) (*sys.RoleIndexResp, error) {
	l := rolelogic.NewRoleIndexLogic(ctx, s.svcCtx)
	return l.RoleIndex(in)
}

func (s *RoleServer) RoleUpdate(ctx context.Context, in *sys.RoleUpdateReq) (*sys.Response, error) {
	l := rolelogic.NewRoleUpdateLogic(ctx, s.svcCtx)
	return l.RoleUpdate(in)
}

func (s *RoleServer) RoleDelete(ctx context.Context, in *sys.RoleDeleteReq) (*sys.Response, error) {
	l := rolelogic.NewRoleDeleteLogic(ctx, s.svcCtx)
	return l.RoleDelete(in)
}

func (s *RoleServer) RoleMenuUpdate(ctx context.Context, in *sys.RoleMenuUpdateReq) (*sys.Response, error) {
	l := rolelogic.NewRoleMenuUpdateLogic(ctx, s.svcCtx)
	return l.RoleMenuUpdate(in)
}
