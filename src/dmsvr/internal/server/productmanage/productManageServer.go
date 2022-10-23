// Code generated by goctl. DO NOT EDIT!
// Source: dm.proto

package server

import (
	"context"

	"github.com/i-Things/things/src/dmsvr/internal/logic/productmanage"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
)

type ProductManageServer struct {
	svcCtx *svc.ServiceContext
	dm.UnimplementedProductManageServer
}

func NewProductManageServer(svcCtx *svc.ServiceContext) *ProductManageServer {
	return &ProductManageServer{
		svcCtx: svcCtx,
	}
}

// 新增产品
func (s *ProductManageServer) ProductInfoCreate(ctx context.Context, in *dm.ProductInfo) (*dm.Response, error) {
	l := productmanagelogic.NewProductInfoCreateLogic(ctx, s.svcCtx)
	return l.ProductInfoCreate(in)
}

// 更新产品
func (s *ProductManageServer) ProductInfoUpdate(ctx context.Context, in *dm.ProductInfo) (*dm.Response, error) {
	l := productmanagelogic.NewProductInfoUpdateLogic(ctx, s.svcCtx)
	return l.ProductInfoUpdate(in)
}

// 删除产品
func (s *ProductManageServer) ProductInfoDelete(ctx context.Context, in *dm.ProductInfoDeleteReq) (*dm.Response, error) {
	l := productmanagelogic.NewProductInfoDeleteLogic(ctx, s.svcCtx)
	return l.ProductInfoDelete(in)
}

// 获取产品信息列表
func (s *ProductManageServer) ProductInfoIndex(ctx context.Context, in *dm.ProductInfoIndexReq) (*dm.ProductInfoIndexResp, error) {
	l := productmanagelogic.NewProductInfoIndexLogic(ctx, s.svcCtx)
	return l.ProductInfoIndex(in)
}

// 获取产品信息详情
func (s *ProductManageServer) ProductInfoRead(ctx context.Context, in *dm.ProductInfoReadReq) (*dm.ProductInfo, error) {
	l := productmanagelogic.NewProductInfoReadLogic(ctx, s.svcCtx)
	return l.ProductInfoRead(in)
}

// 更新产品物模型
func (s *ProductManageServer) ProductSchemaUpdate(ctx context.Context, in *dm.ProductSchemaUpdateReq) (*dm.Response, error) {
	l := productmanagelogic.NewProductSchemaUpdateLogic(ctx, s.svcCtx)
	return l.ProductSchemaUpdate(in)
}

// 新增产品
func (s *ProductManageServer) ProductSchemaCreate(ctx context.Context, in *dm.ProductSchemaCreateReq) (*dm.Response, error) {
	l := productmanagelogic.NewProductSchemaCreateLogic(ctx, s.svcCtx)
	return l.ProductSchemaCreate(in)
}

// 删除产品
func (s *ProductManageServer) ProductSchemaDelete(ctx context.Context, in *dm.ProductSchemaDeleteReq) (*dm.Response, error) {
	l := productmanagelogic.NewProductSchemaDeleteLogic(ctx, s.svcCtx)
	return l.ProductSchemaDelete(in)
}

// 获取产品信息列表
func (s *ProductManageServer) ProductSchemaIndex(ctx context.Context, in *dm.ProductSchemaIndexReq) (*dm.ProductSchemaIndexResp, error) {
	l := productmanagelogic.NewProductSchemaIndexLogic(ctx, s.svcCtx)
	return l.ProductSchemaIndex(in)
}

// 删除产品
func (s *ProductManageServer) ProductSchemaTslImport(ctx context.Context, in *dm.ProductSchemaTslImportReq) (*dm.Response, error) {
	l := productmanagelogic.NewProductSchemaTslImportLogic(ctx, s.svcCtx)
	return l.ProductSchemaTslImport(in)
}

// 获取产品信息列表
func (s *ProductManageServer) ProductSchemaTslRead(ctx context.Context, in *dm.ProductSchemaTslReadReq) (*dm.ProductSchemaTslReadResp, error) {
	l := productmanagelogic.NewProductSchemaTslReadLogic(ctx, s.svcCtx)
	return l.ProductSchemaTslRead(in)
}
