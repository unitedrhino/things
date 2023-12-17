// Code generated by goctl. DO NOT EDIT.
// Source: sys.proto

package appmanage

import (
	"context"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	ApiDeleteReq                  = sys.ApiDeleteReq
	ApiInfo                       = sys.ApiInfo
	ApiInfoIndexReq               = sys.ApiInfoIndexReq
	ApiInfoIndexResp              = sys.ApiInfoIndexResp
	AppInfo                       = sys.AppInfo
	AppInfoIndexReq               = sys.AppInfoIndexReq
	AppInfoIndexResp              = sys.AppInfoIndexResp
	AreaInfo                      = sys.AreaInfo
	AreaInfoDeleteReq             = sys.AreaInfoDeleteReq
	AreaInfoIndexReq              = sys.AreaInfoIndexReq
	AreaInfoIndexResp             = sys.AreaInfoIndexResp
	AreaInfoReadReq               = sys.AreaInfoReadReq
	AreaInfoTreeReq               = sys.AreaInfoTreeReq
	AreaInfoTreeResp              = sys.AreaInfoTreeResp
	AuthApiInfo                   = sys.AuthApiInfo
	ConfigResp                    = sys.ConfigResp
	DateRange                     = sys.DateRange
	JwtToken                      = sys.JwtToken
	LoginLogCreateReq             = sys.LoginLogCreateReq
	LoginLogIndexReq              = sys.LoginLogIndexReq
	LoginLogIndexResp             = sys.LoginLogIndexResp
	LoginLogInfo                  = sys.LoginLogInfo
	Map                           = sys.Map
	MenuInfo                      = sys.MenuInfo
	MenuInfoIndexReq              = sys.MenuInfoIndexReq
	MenuInfoIndexResp             = sys.MenuInfoIndexResp
	OperLogCreateReq              = sys.OperLogCreateReq
	OperLogIndexReq               = sys.OperLogIndexReq
	OperLogIndexResp              = sys.OperLogIndexResp
	OperLogInfo                   = sys.OperLogInfo
	PageInfo                      = sys.PageInfo
	PageInfo_OrderBy              = sys.PageInfo_OrderBy
	Point                         = sys.Point
	ProjectInfo                   = sys.ProjectInfo
	ProjectInfoDeleteReq          = sys.ProjectInfoDeleteReq
	ProjectInfoIndexReq           = sys.ProjectInfoIndexReq
	ProjectInfoIndexResp          = sys.ProjectInfoIndexResp
	ProjectInfoReadReq            = sys.ProjectInfoReadReq
	ReqWithID                     = sys.ReqWithID
	ReqWithIDCode                 = sys.ReqWithIDCode
	Response                      = sys.Response
	RoleApiAuthReq                = sys.RoleApiAuthReq
	RoleApiIndexReq               = sys.RoleApiIndexReq
	RoleApiIndexResp              = sys.RoleApiIndexResp
	RoleApiMultiUpdateReq         = sys.RoleApiMultiUpdateReq
	RoleAppIndexReq               = sys.RoleAppIndexReq
	RoleAppIndexResp              = sys.RoleAppIndexResp
	RoleAppMultiUpdateReq         = sys.RoleAppMultiUpdateReq
	RoleAppUpdateReq              = sys.RoleAppUpdateReq
	RoleInfo                      = sys.RoleInfo
	RoleInfoIndexReq              = sys.RoleInfoIndexReq
	RoleInfoIndexResp             = sys.RoleInfoIndexResp
	RoleMenuIndexReq              = sys.RoleMenuIndexReq
	RoleMenuIndexResp             = sys.RoleMenuIndexResp
	RoleMenuMultiUpdateReq        = sys.RoleMenuMultiUpdateReq
	TenantAppIndexReq             = sys.TenantAppIndexReq
	TenantAppIndexResp            = sys.TenantAppIndexResp
	TenantAppMultiUpdateReq       = sys.TenantAppMultiUpdateReq
	TenantInfo                    = sys.TenantInfo
	TenantInfoIndexReq            = sys.TenantInfoIndexReq
	TenantInfoIndexResp           = sys.TenantInfoIndexResp
	UserAuthArea                  = sys.UserAuthArea
	UserAuthAreaIndexReq          = sys.UserAuthAreaIndexReq
	UserAuthAreaIndexResp         = sys.UserAuthAreaIndexResp
	UserAuthAreaMultiUpdateReq    = sys.UserAuthAreaMultiUpdateReq
	UserAuthProject               = sys.UserAuthProject
	UserAuthProjectIndexReq       = sys.UserAuthProjectIndexReq
	UserAuthProjectIndexResp      = sys.UserAuthProjectIndexResp
	UserAuthProjectMultiUpdateReq = sys.UserAuthProjectMultiUpdateReq
	UserCheckTokenReq             = sys.UserCheckTokenReq
	UserCheckTokenResp            = sys.UserCheckTokenResp
	UserCreateResp                = sys.UserCreateResp
	UserInfo                      = sys.UserInfo
	UserInfoDeleteReq             = sys.UserInfoDeleteReq
	UserInfoIndexReq              = sys.UserInfoIndexReq
	UserInfoIndexResp             = sys.UserInfoIndexResp
	UserInfoReadReq               = sys.UserInfoReadReq
	UserLoginReq                  = sys.UserLoginReq
	UserLoginResp                 = sys.UserLoginResp
	UserRegister1Req              = sys.UserRegister1Req
	UserRegister1Resp             = sys.UserRegister1Resp
	UserRegister2Req              = sys.UserRegister2Req
	UserRoleIndexReq              = sys.UserRoleIndexReq
	UserRoleIndexResp             = sys.UserRoleIndexResp
	UserRoleMultiUpdateReq        = sys.UserRoleMultiUpdateReq

	AppManage interface {
		AppInfoCreate(ctx context.Context, in *AppInfo, opts ...grpc.CallOption) (*Response, error)
		AppInfoIndex(ctx context.Context, in *AppInfoIndexReq, opts ...grpc.CallOption) (*AppInfoIndexResp, error)
		AppInfoUpdate(ctx context.Context, in *AppInfo, opts ...grpc.CallOption) (*Response, error)
		AppInfoDelete(ctx context.Context, in *ReqWithIDCode, opts ...grpc.CallOption) (*Response, error)
		AppInfoRead(ctx context.Context, in *ReqWithIDCode, opts ...grpc.CallOption) (*AppInfo, error)
	}

	defaultAppManage struct {
		cli zrpc.Client
	}

	directAppManage struct {
		svcCtx *svc.ServiceContext
		svr    sys.AppManageServer
	}
)

func NewAppManage(cli zrpc.Client) AppManage {
	return &defaultAppManage{
		cli: cli,
	}
}

func NewDirectAppManage(svcCtx *svc.ServiceContext, svr sys.AppManageServer) AppManage {
	return &directAppManage{
		svr:    svr,
		svcCtx: svcCtx,
	}
}

func (m *defaultAppManage) AppInfoCreate(ctx context.Context, in *AppInfo, opts ...grpc.CallOption) (*Response, error) {
	client := sys.NewAppManageClient(m.cli.Conn())
	return client.AppInfoCreate(ctx, in, opts...)
}

func (d *directAppManage) AppInfoCreate(ctx context.Context, in *AppInfo, opts ...grpc.CallOption) (*Response, error) {
	return d.svr.AppInfoCreate(ctx, in)
}

func (m *defaultAppManage) AppInfoIndex(ctx context.Context, in *AppInfoIndexReq, opts ...grpc.CallOption) (*AppInfoIndexResp, error) {
	client := sys.NewAppManageClient(m.cli.Conn())
	return client.AppInfoIndex(ctx, in, opts...)
}

func (d *directAppManage) AppInfoIndex(ctx context.Context, in *AppInfoIndexReq, opts ...grpc.CallOption) (*AppInfoIndexResp, error) {
	return d.svr.AppInfoIndex(ctx, in)
}

func (m *defaultAppManage) AppInfoUpdate(ctx context.Context, in *AppInfo, opts ...grpc.CallOption) (*Response, error) {
	client := sys.NewAppManageClient(m.cli.Conn())
	return client.AppInfoUpdate(ctx, in, opts...)
}

func (d *directAppManage) AppInfoUpdate(ctx context.Context, in *AppInfo, opts ...grpc.CallOption) (*Response, error) {
	return d.svr.AppInfoUpdate(ctx, in)
}

func (m *defaultAppManage) AppInfoDelete(ctx context.Context, in *ReqWithIDCode, opts ...grpc.CallOption) (*Response, error) {
	client := sys.NewAppManageClient(m.cli.Conn())
	return client.AppInfoDelete(ctx, in, opts...)
}

func (d *directAppManage) AppInfoDelete(ctx context.Context, in *ReqWithIDCode, opts ...grpc.CallOption) (*Response, error) {
	return d.svr.AppInfoDelete(ctx, in)
}

func (m *defaultAppManage) AppInfoRead(ctx context.Context, in *ReqWithIDCode, opts ...grpc.CallOption) (*AppInfo, error) {
	client := sys.NewAppManageClient(m.cli.Conn())
	return client.AppInfoRead(ctx, in, opts...)
}

func (d *directAppManage) AppInfoRead(ctx context.Context, in *ReqWithIDCode, opts ...grpc.CallOption) (*AppInfo, error) {
	return d.svr.AppInfoRead(ctx, in)
}