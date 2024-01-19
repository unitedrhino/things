package ctxs

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net/http"
)

type UserCtx struct {
	IsOpen     bool //是否开放认证用户
	AppCode    string
	TenantCode string //租户Code
	ProjectID  int64  `json:",string"`
	IsAdmin    bool   //是否是超级管理员
	UserID     int64  `json:",string"` //用户id（开放认证用户值为0）
	RoleID     int64  //用户使用的角色（开放认证用户值为0）
	IsAllData  bool   //是否所有数据权限（开放认证用户值为true）
	IP         string //用户的ip地址
	Os         string //操作系统
	InnerCtx
}

type InnerCtx struct {
	AllProject bool
	AllArea    bool //内部使用,不限制区域
	AllTenant  bool //所有租户的权限
}

func NotLoginedInit(r *http.Request) *http.Request {
	strIP, _ := utils.GetIP(r)
	appCode := r.Header.Get(UserAppCodeKey)
	if appCode == "" {
		appCode = def.AppCore
	}
	tenantCode := r.Header.Get(UserTenantCodeKey)
	if tenantCode == "" {
		tenantCode = def.TenantCodeDefault
	}
	c := context.WithValue(r.Context(), UserInfoKey, &UserCtx{
		AppCode:    appCode,
		TenantCode: tenantCode,
		IP:         strIP,
		Os:         r.Header.Get("User-Agent"),
	})
	return r.WithContext(c)
}

func SetUserCtx(ctx context.Context, userCtx *UserCtx) context.Context {
	info, _ := json.Marshal(userCtx)
	ctx = metadata.NewOutgoingContext(context.Background(), metadata.Pairs(
		UserInfoKey, string(info),
	))
	return context.WithValue(ctx, UserInfoKey, userCtx)
}
func SetInnerCtx(ctx context.Context, inner InnerCtx) context.Context {
	uc := GetUserCtx(ctx)
	if uc == nil {
		return ctx
	}
	uc.InnerCtx = inner
	return SetUserCtx(ctx, uc)
}

func GetInnerCtx(ctx context.Context) InnerCtx {
	uc := GetUserCtx(ctx)
	if uc == nil {
		return InnerCtx{}
	}
	return uc.InnerCtx
}

// 使用该函数前必须传了UserCtx
func GetUserCtx(ctx context.Context) *UserCtx {
	val, ok := ctx.Value(UserInfoKey).(*UserCtx)
	if !ok { //这里线上不能获取不到
		return nil
	}
	return val
}

func GrpcInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	ctx = func() context.Context {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return ctx
		}
		info := md[UserInfoKey]
		if len(info) == 0 {
			return ctx
		}
		var val UserCtx
		if err := json.Unmarshal([]byte(info[0]), &val); err != nil {
			return ctx
		}
		return SetUserCtx(ctx, &val)
	}()
	resp, err := handler(ctx, req)
	return resp, err
}

func IsRoot(ctx context.Context) error {
	uc := GetUserCtx(ctx)
	if uc == nil || uc.TenantCode != def.TenantCodeDefault {
		return errors.Permissions.AddDetailf("需要主租户才能操作")
	}
	return nil
}

// 使用该函数前必须传了UserCtx
func GetUserCtxOrNil(ctx context.Context) *UserCtx {
	val, ok := ctx.Value(UserInfoKey).(*UserCtx)
	if !ok { //这里线上不能获取不到
		return nil
	}
	return val
}

type MetadataCtx = map[string][]string

func SetMetaCtx(ctx context.Context, maps MetadataCtx) context.Context {
	return context.WithValue(ctx, MetadataKey, maps)
}
func GetMetaCtx(ctx context.Context) MetadataCtx {
	val, ok := ctx.Value(MetadataKey).(MetadataCtx)
	if !ok {
		return nil
	}
	return val
}

func GetMetaVal(ctx context.Context, field string) []string {
	mdCtx := GetMetaCtx(ctx)
	if val, ok := mdCtx[field]; !ok {
		return nil
	} else {
		return val
	}
}

// 指定项目id（企业版功能）
func SetMetaProjectID(ctx context.Context, projectID int64) {
	mc := GetMetaCtx(ctx)
	projectIDStr := utils.ToString(projectID)
	mc[string(MetaFieldProjectID)] = []string{projectIDStr}
}

// 获取meta里的项目ID（企业版功能）
func ClearMetaProjectID(ctx context.Context) {
	mc := GetMetaCtx(ctx)
	delete(mc, string(MetaFieldProjectID))
}
