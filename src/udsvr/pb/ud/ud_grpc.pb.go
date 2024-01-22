// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.19.4
// source: proto/ud.proto

package ud

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Rule_SceneInfoCreate_FullMethodName      = "/ud.rule/sceneInfoCreate"
	Rule_SceneInfoUpdate_FullMethodName      = "/ud.rule/sceneInfoUpdate"
	Rule_SceneInfoDelete_FullMethodName      = "/ud.rule/sceneInfoDelete"
	Rule_SceneInfoIndex_FullMethodName       = "/ud.rule/sceneInfoIndex"
	Rule_SceneInfoRead_FullMethodName        = "/ud.rule/sceneInfoRead"
	Rule_SceneManuallyTrigger_FullMethodName = "/ud.rule/sceneManuallyTrigger"
)

// RuleClient is the client API for Rule service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RuleClient interface {
	// 场景
	SceneInfoCreate(ctx context.Context, in *SceneInfo, opts ...grpc.CallOption) (*WithID, error)
	SceneInfoUpdate(ctx context.Context, in *SceneInfo, opts ...grpc.CallOption) (*Empty, error)
	SceneInfoDelete(ctx context.Context, in *WithID, opts ...grpc.CallOption) (*Empty, error)
	SceneInfoIndex(ctx context.Context, in *SceneInfoIndexReq, opts ...grpc.CallOption) (*SceneInfoIndexResp, error)
	SceneInfoRead(ctx context.Context, in *WithID, opts ...grpc.CallOption) (*SceneInfo, error)
	SceneManuallyTrigger(ctx context.Context, in *WithID, opts ...grpc.CallOption) (*Empty, error)
}

type ruleClient struct {
	cc grpc.ClientConnInterface
}

func NewRuleClient(cc grpc.ClientConnInterface) RuleClient {
	return &ruleClient{cc}
}

func (c *ruleClient) SceneInfoCreate(ctx context.Context, in *SceneInfo, opts ...grpc.CallOption) (*WithID, error) {
	out := new(WithID)
	err := c.cc.Invoke(ctx, Rule_SceneInfoCreate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ruleClient) SceneInfoUpdate(ctx context.Context, in *SceneInfo, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, Rule_SceneInfoUpdate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ruleClient) SceneInfoDelete(ctx context.Context, in *WithID, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, Rule_SceneInfoDelete_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ruleClient) SceneInfoIndex(ctx context.Context, in *SceneInfoIndexReq, opts ...grpc.CallOption) (*SceneInfoIndexResp, error) {
	out := new(SceneInfoIndexResp)
	err := c.cc.Invoke(ctx, Rule_SceneInfoIndex_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ruleClient) SceneInfoRead(ctx context.Context, in *WithID, opts ...grpc.CallOption) (*SceneInfo, error) {
	out := new(SceneInfo)
	err := c.cc.Invoke(ctx, Rule_SceneInfoRead_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ruleClient) SceneManuallyTrigger(ctx context.Context, in *WithID, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, Rule_SceneManuallyTrigger_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RuleServer is the server API for Rule service.
// All implementations must embed UnimplementedRuleServer
// for forward compatibility
type RuleServer interface {
	// 场景
	SceneInfoCreate(context.Context, *SceneInfo) (*WithID, error)
	SceneInfoUpdate(context.Context, *SceneInfo) (*Empty, error)
	SceneInfoDelete(context.Context, *WithID) (*Empty, error)
	SceneInfoIndex(context.Context, *SceneInfoIndexReq) (*SceneInfoIndexResp, error)
	SceneInfoRead(context.Context, *WithID) (*SceneInfo, error)
	SceneManuallyTrigger(context.Context, *WithID) (*Empty, error)
	mustEmbedUnimplementedRuleServer()
}

// UnimplementedRuleServer must be embedded to have forward compatible implementations.
type UnimplementedRuleServer struct {
}

func (UnimplementedRuleServer) SceneInfoCreate(context.Context, *SceneInfo) (*WithID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SceneInfoCreate not implemented")
}
func (UnimplementedRuleServer) SceneInfoUpdate(context.Context, *SceneInfo) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SceneInfoUpdate not implemented")
}
func (UnimplementedRuleServer) SceneInfoDelete(context.Context, *WithID) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SceneInfoDelete not implemented")
}
func (UnimplementedRuleServer) SceneInfoIndex(context.Context, *SceneInfoIndexReq) (*SceneInfoIndexResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SceneInfoIndex not implemented")
}
func (UnimplementedRuleServer) SceneInfoRead(context.Context, *WithID) (*SceneInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SceneInfoRead not implemented")
}
func (UnimplementedRuleServer) SceneManuallyTrigger(context.Context, *WithID) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SceneManuallyTrigger not implemented")
}
func (UnimplementedRuleServer) mustEmbedUnimplementedRuleServer() {}

// UnsafeRuleServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RuleServer will
// result in compilation errors.
type UnsafeRuleServer interface {
	mustEmbedUnimplementedRuleServer()
}

func RegisterRuleServer(s grpc.ServiceRegistrar, srv RuleServer) {
	s.RegisterService(&Rule_ServiceDesc, srv)
}

func _Rule_SceneInfoCreate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SceneInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuleServer).SceneInfoCreate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Rule_SceneInfoCreate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuleServer).SceneInfoCreate(ctx, req.(*SceneInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rule_SceneInfoUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SceneInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuleServer).SceneInfoUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Rule_SceneInfoUpdate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuleServer).SceneInfoUpdate(ctx, req.(*SceneInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rule_SceneInfoDelete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WithID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuleServer).SceneInfoDelete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Rule_SceneInfoDelete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuleServer).SceneInfoDelete(ctx, req.(*WithID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rule_SceneInfoIndex_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SceneInfoIndexReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuleServer).SceneInfoIndex(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Rule_SceneInfoIndex_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuleServer).SceneInfoIndex(ctx, req.(*SceneInfoIndexReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rule_SceneInfoRead_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WithID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuleServer).SceneInfoRead(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Rule_SceneInfoRead_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuleServer).SceneInfoRead(ctx, req.(*WithID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rule_SceneManuallyTrigger_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WithID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuleServer).SceneManuallyTrigger(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Rule_SceneManuallyTrigger_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuleServer).SceneManuallyTrigger(ctx, req.(*WithID))
	}
	return interceptor(ctx, in, info, handler)
}

// Rule_ServiceDesc is the grpc.ServiceDesc for Rule service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Rule_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ud.rule",
	HandlerType: (*RuleServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "sceneInfoCreate",
			Handler:    _Rule_SceneInfoCreate_Handler,
		},
		{
			MethodName: "sceneInfoUpdate",
			Handler:    _Rule_SceneInfoUpdate_Handler,
		},
		{
			MethodName: "sceneInfoDelete",
			Handler:    _Rule_SceneInfoDelete_Handler,
		},
		{
			MethodName: "sceneInfoIndex",
			Handler:    _Rule_SceneInfoIndex_Handler,
		},
		{
			MethodName: "sceneInfoRead",
			Handler:    _Rule_SceneInfoRead_Handler,
		},
		{
			MethodName: "sceneManuallyTrigger",
			Handler:    _Rule_SceneManuallyTrigger_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/ud.proto",
}

const (
	Ops_OpsWorkOrderCreate_FullMethodName = "/ud.ops/opsWorkOrderCreate"
	Ops_OpsWorkOrderUpdate_FullMethodName = "/ud.ops/opsWorkOrderUpdate"
	Ops_OpsWorkOrderIndex_FullMethodName  = "/ud.ops/opsWorkOrderIndex"
)

// OpsClient is the client API for Ops service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OpsClient interface {
	// 维护工单  Work Order
	OpsWorkOrderCreate(ctx context.Context, in *OpsWorkOrder, opts ...grpc.CallOption) (*WithID, error)
	OpsWorkOrderUpdate(ctx context.Context, in *OpsWorkOrder, opts ...grpc.CallOption) (*Empty, error)
	OpsWorkOrderIndex(ctx context.Context, in *OpsWorkOrderIndexReq, opts ...grpc.CallOption) (*OpsWorkOrderIndexResp, error)
}

type opsClient struct {
	cc grpc.ClientConnInterface
}

func NewOpsClient(cc grpc.ClientConnInterface) OpsClient {
	return &opsClient{cc}
}

func (c *opsClient) OpsWorkOrderCreate(ctx context.Context, in *OpsWorkOrder, opts ...grpc.CallOption) (*WithID, error) {
	out := new(WithID)
	err := c.cc.Invoke(ctx, Ops_OpsWorkOrderCreate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *opsClient) OpsWorkOrderUpdate(ctx context.Context, in *OpsWorkOrder, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, Ops_OpsWorkOrderUpdate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *opsClient) OpsWorkOrderIndex(ctx context.Context, in *OpsWorkOrderIndexReq, opts ...grpc.CallOption) (*OpsWorkOrderIndexResp, error) {
	out := new(OpsWorkOrderIndexResp)
	err := c.cc.Invoke(ctx, Ops_OpsWorkOrderIndex_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OpsServer is the server API for Ops service.
// All implementations must embed UnimplementedOpsServer
// for forward compatibility
type OpsServer interface {
	// 维护工单  Work Order
	OpsWorkOrderCreate(context.Context, *OpsWorkOrder) (*WithID, error)
	OpsWorkOrderUpdate(context.Context, *OpsWorkOrder) (*Empty, error)
	OpsWorkOrderIndex(context.Context, *OpsWorkOrderIndexReq) (*OpsWorkOrderIndexResp, error)
	mustEmbedUnimplementedOpsServer()
}

// UnimplementedOpsServer must be embedded to have forward compatible implementations.
type UnimplementedOpsServer struct {
}

func (UnimplementedOpsServer) OpsWorkOrderCreate(context.Context, *OpsWorkOrder) (*WithID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OpsWorkOrderCreate not implemented")
}
func (UnimplementedOpsServer) OpsWorkOrderUpdate(context.Context, *OpsWorkOrder) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OpsWorkOrderUpdate not implemented")
}
func (UnimplementedOpsServer) OpsWorkOrderIndex(context.Context, *OpsWorkOrderIndexReq) (*OpsWorkOrderIndexResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OpsWorkOrderIndex not implemented")
}
func (UnimplementedOpsServer) mustEmbedUnimplementedOpsServer() {}

// UnsafeOpsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OpsServer will
// result in compilation errors.
type UnsafeOpsServer interface {
	mustEmbedUnimplementedOpsServer()
}

func RegisterOpsServer(s grpc.ServiceRegistrar, srv OpsServer) {
	s.RegisterService(&Ops_ServiceDesc, srv)
}

func _Ops_OpsWorkOrderCreate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OpsWorkOrder)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OpsServer).OpsWorkOrderCreate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Ops_OpsWorkOrderCreate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsServer).OpsWorkOrderCreate(ctx, req.(*OpsWorkOrder))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ops_OpsWorkOrderUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OpsWorkOrder)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OpsServer).OpsWorkOrderUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Ops_OpsWorkOrderUpdate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsServer).OpsWorkOrderUpdate(ctx, req.(*OpsWorkOrder))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ops_OpsWorkOrderIndex_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OpsWorkOrderIndexReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OpsServer).OpsWorkOrderIndex(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Ops_OpsWorkOrderIndex_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsServer).OpsWorkOrderIndex(ctx, req.(*OpsWorkOrderIndexReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Ops_ServiceDesc is the grpc.ServiceDesc for Ops service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Ops_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ud.ops",
	HandlerType: (*OpsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "opsWorkOrderCreate",
			Handler:    _Ops_OpsWorkOrderCreate_Handler,
		},
		{
			MethodName: "opsWorkOrderUpdate",
			Handler:    _Ops_OpsWorkOrderUpdate_Handler,
		},
		{
			MethodName: "opsWorkOrderIndex",
			Handler:    _Ops_OpsWorkOrderIndex_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/ud.proto",
}
