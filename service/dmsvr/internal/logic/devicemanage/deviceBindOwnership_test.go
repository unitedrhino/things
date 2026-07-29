package devicemanagelogic

import (
	"context"
	"testing"

	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"google.golang.org/grpc"
)

type bindProjectReaderFunc func(context.Context, *sys.ProjectWithID, ...grpc.CallOption) (*sys.ProjectInfo, error)

func (f bindProjectReaderFunc) ProjectInfoRead(ctx context.Context, in *sys.ProjectWithID, opts ...grpc.CallOption) (*sys.ProjectInfo, error) {
	return f(ctx, in, opts...)
}

type bindUserReaderFunc func(context.Context, *sys.UserInfoReadReq, ...grpc.CallOption) (*sys.UserInfo, error)

func (f bindUserReaderFunc) UserInfoRead(ctx context.Context, in *sys.UserInfoReadReq, opts ...grpc.CallOption) (*sys.UserInfo, error) {
	return f(ctx, in, opts...)
}

func TestClassifyDeviceBindOwnershipAllowsDeletedProjectAdmin(t *testing.T) {
	got, err := classifyDeviceBindOwnership(context.Background(),
		bindProjectReaderFunc(func(context.Context, *sys.ProjectWithID, ...grpc.CallOption) (*sys.ProjectInfo, error) {
			return &sys.ProjectInfo{ProjectID: 9001, AdminUserID: 8001}, nil
		}),
		bindUserReaderFunc(func(context.Context, *sys.UserInfoReadReq, ...grpc.CallOption) (*sys.UserInfo, error) {
			return nil, errors.NotFind
		}),
		def.TenantCodeDefault,
		9001,
		&ctxs.UserCtx{TenantCode: def.TenantCodeDefault, ProjectID: 7001},
		0,
	)
	if err != nil {
		t.Fatalf("classifyDeviceBindOwnership error: %v", err)
	}
	if got != deviceBindOwnershipStale {
		t.Fatalf("ownership = %v, want stale", got)
	}
}

func TestClassifyDeviceBindOwnershipBlocksActiveOtherProject(t *testing.T) {
	got, err := classifyDeviceBindOwnership(context.Background(),
		bindProjectReaderFunc(func(context.Context, *sys.ProjectWithID, ...grpc.CallOption) (*sys.ProjectInfo, error) {
			return &sys.ProjectInfo{ProjectID: 9001, AdminUserID: 8001}, nil
		}),
		bindUserReaderFunc(func(context.Context, *sys.UserInfoReadReq, ...grpc.CallOption) (*sys.UserInfo, error) {
			return &sys.UserInfo{UserID: 8001}, nil
		}),
		def.TenantCodeDefault,
		9001,
		&ctxs.UserCtx{TenantCode: def.TenantCodeDefault, ProjectID: 7001},
		0,
	)
	if err != nil {
		t.Fatalf("classifyDeviceBindOwnership error: %v", err)
	}
	if got != deviceBindOwnershipBlocked {
		t.Fatalf("ownership = %v, want blocked", got)
	}
}

func TestClassifyDeviceBindOwnershipDetectsCurrentUserOtherProject(t *testing.T) {
	got, err := classifyDeviceBindOwnership(context.Background(),
		bindProjectReaderFunc(func(context.Context, *sys.ProjectWithID, ...grpc.CallOption) (*sys.ProjectInfo, error) {
			return &sys.ProjectInfo{ProjectID: 9001, AdminUserID: 7001}, nil
		}),
		bindUserReaderFunc(func(context.Context, *sys.UserInfoReadReq, ...grpc.CallOption) (*sys.UserInfo, error) {
			t.Fatal("user reader should not be called when current user owns project")
			return nil, nil
		}),
		def.TenantCodeDefault,
		9001,
		&ctxs.UserCtx{TenantCode: def.TenantCodeDefault, ProjectID: 7002, UserID: 7001},
		0,
	)
	if err != nil {
		t.Fatalf("classifyDeviceBindOwnership error: %v", err)
	}
	if got != deviceBindOwnershipCurrentUserBound {
		t.Fatalf("ownership = %v, want current user bound", got)
	}
}

func TestClassifyDeviceBindOwnershipAllowsMissingProject(t *testing.T) {
	got, err := classifyDeviceBindOwnership(context.Background(),
		bindProjectReaderFunc(func(context.Context, *sys.ProjectWithID, ...grpc.CallOption) (*sys.ProjectInfo, error) {
			return nil, errors.NotFind
		}),
		bindUserReaderFunc(func(context.Context, *sys.UserInfoReadReq, ...grpc.CallOption) (*sys.UserInfo, error) {
			t.Fatal("user reader should not be called when project is missing")
			return nil, nil
		}),
		def.TenantCodeDefault,
		9001,
		&ctxs.UserCtx{TenantCode: def.TenantCodeDefault, ProjectID: 7001},
		0,
	)
	if err != nil {
		t.Fatalf("classifyDeviceBindOwnership error: %v", err)
	}
	if got != deviceBindOwnershipStale {
		t.Fatalf("ownership = %v, want stale", got)
	}
}

func TestClassifyDeviceBindOwnershipReturnsUserReadError(t *testing.T) {
	wantErr := errors.Database.AddMsg("user lookup failed")
	got, err := classifyDeviceBindOwnership(context.Background(),
		bindProjectReaderFunc(func(context.Context, *sys.ProjectWithID, ...grpc.CallOption) (*sys.ProjectInfo, error) {
			return &sys.ProjectInfo{ProjectID: 9001, AdminUserID: 8001}, nil
		}),
		bindUserReaderFunc(func(context.Context, *sys.UserInfoReadReq, ...grpc.CallOption) (*sys.UserInfo, error) {
			return nil, wantErr
		}),
		def.TenantCodeDefault,
		9001,
		&ctxs.UserCtx{TenantCode: def.TenantCodeDefault, ProjectID: 7001},
		0,
	)
	if err == nil {
		t.Fatal("classifyDeviceBindOwnership error = nil, want user read error")
	}
	if !errors.Cmp(err, wantErr) {
		t.Fatalf("classifyDeviceBindOwnership error = %v, want %v", err, wantErr)
	}
	if got != deviceBindOwnershipBlocked {
		t.Fatalf("ownership = %v, want blocked", got)
	}
}
