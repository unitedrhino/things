package cache

import (
	"context"
	schema "gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"gitee.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/dgraph-io/ristretto"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

const (
	expirtTime = time.Hour
)

type SchemaRepo struct {
	cache *ristretto.Cache
}

func NewSchemaRepo() schema.Repo {
	cache, _ := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	return &SchemaRepo{
		cache: cache,
	}
}

func (s SchemaRepo) TslImport(ctx context.Context, productID string, schemaInfo *schema.Model) error {
	db := relationDB.NewProductSchemaRepo(ctx)
	//todo 后续需要修改为事务处理
	err := s.Delete(ctx, productID)
	if err != nil {
		logx.WithContext(ctx).Errorf("%s.Delete err:%v", utils.FuncName(), err)
		return errors.Database
	}
	for _, item := range schemaInfo.Property {
		err = db.Insert(ctx, relationDB.ToPropertyPo(productID, item))
		if err != nil {
			return err
		}
	}
	for _, item := range schemaInfo.Event {
		err = db.Insert(ctx, relationDB.ToEventPo(productID, item))
		if err != nil {
			return err
		}
	}
	for _, item := range schemaInfo.Action {
		err = db.Insert(ctx, relationDB.ToActionPo(productID, item))
		if err != nil {
			return err
		}
	}
	s.cache.SetWithTTL(productID, schemaInfo, 1, expirtTime)
	return err
}

//func (s SchemaRepo) TslRead(ctx context.Context, productID string) (*schema.Model, error) {
//	temp, ok := s.cache.Get(productID)
//	if ok {
//		return temp.(*schema.Model), nil
//	}
//	db := relationDB.NewProductSchemaRepo(ctx)
//	dbSchemas, err := db.FindByFilter(ctx, relationDB.ProductSchemaFilter{ProductID: productID}, nil)
//	if err != nil {
//		return nil, err
//	}
//	schemaModel := relationDB.ToSchemaDo(productID, dbSchemas)
//	s.cache.SetWithTTL(productID, schemaModel, 1, expirtTime)
//	return schemaModel, nil
//}

func (s SchemaRepo) GetSchemaModel(ctx context.Context, productID string) (*schema.Model, error) {
	temp, ok := s.cache.Get(productID)
	if ok {
		return temp.(*schema.Model), nil
	}
	db := relationDB.NewProductSchemaRepo(ctx)
	dbSchemas, err := db.FindByFilter(ctx, relationDB.ProductSchemaFilter{ProductID: productID}, nil)
	if err != nil {
		return nil, err
	}
	schemaModel := relationDB.ToSchemaDo(productID, dbSchemas)
	s.cache.SetWithTTL(productID, schemaModel, 1, expirtTime)
	return schemaModel, nil
}

func (s SchemaRepo) ClearCache(ctx context.Context, productID string) error {
	s.cache.Del(productID)
	return nil
}

func (s SchemaRepo) Delete(ctx context.Context, productID string) error {
	s.cache.Del(productID)
	db := relationDB.NewProductSchemaRepo(ctx)
	err := db.DeleteByFilter(ctx, relationDB.ProductSchemaFilter{ProductID: productID})
	if err != nil {
		logx.WithContext(ctx).Errorf("%s.Delete err:%v", utils.FuncName(), err)
		return errors.Database
	}
	return err
}
