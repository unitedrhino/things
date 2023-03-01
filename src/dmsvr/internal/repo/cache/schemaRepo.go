package cache

import (
	"context"
	"github.com/dgraph-io/ristretto"
	schema "github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

const (
	expirtTime = time.Hour
)

type SchemaRepo struct {
	db    mysql.DmProductSchemaModel
	cache *ristretto.Cache
}

func NewSchemaRepo(t mysql.DmProductSchemaModel) schema.Repo {
	cache, _ := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	return &SchemaRepo{
		db:    t,
		cache: cache,
	}
}

func (s SchemaRepo) TslImport(ctx context.Context, productID string, schemaInfo *schema.Model) error {
	//todo 后续需要修改为事务处理
	err := s.Delete(ctx, productID)
	if err != nil {
		logx.WithContext(ctx).Errorf("%s.Delete err:%v", utils.FuncName(), err)
		return errors.Database
	}
	for _, item := range schemaInfo.Property {
		_, err = s.db.Insert(ctx, mysql.ToPropertyPo(productID, item))
		if err != nil {
			return err
		}
	}
	for _, item := range schemaInfo.Event {
		_, err = s.db.Insert(ctx, mysql.ToEventPo(productID, item))
		if err != nil {
			return err
		}
	}
	for _, item := range schemaInfo.Action {
		_, err = s.db.Insert(ctx, mysql.ToActionPo(productID, item))
		if err != nil {
			return err
		}
	}
	s.cache.SetWithTTL(productID, schemaInfo, 1, expirtTime)
	return err
}
func (s SchemaRepo) TslRead(ctx context.Context, productID string) (*schema.Model, error) {
	temp, ok := s.cache.Get(productID)
	if ok {
		return temp.(*schema.Model), nil
	}
	dbSchemas, err := s.db.FindByFilter(ctx, mysql.ProductSchemaFilter{ProductID: productID}, nil)
	if err != nil {
		return nil, err
	}
	schemaModel := mysql.ToSchemaDo(productID, dbSchemas)
	s.cache.SetWithTTL(productID, schemaModel, 1, expirtTime)
	return schemaModel, nil
}

func (s SchemaRepo) GetSchemaModel(ctx context.Context, productID string) (*schema.Model, error) {
	temp, ok := s.cache.Get(productID)
	if ok {
		return temp.(*schema.Model), nil
	}
	dbSchemas, err := s.db.FindByFilter(ctx, mysql.ProductSchemaFilter{ProductID: productID}, nil)
	if err != nil {
		return nil, err
	}
	schemaModel := mysql.ToSchemaDo(productID, dbSchemas)
	s.cache.SetWithTTL(productID, schemaModel, 1, expirtTime)
	return schemaModel, nil
}

func (s SchemaRepo) ClearCache(ctx context.Context, productID string) error {
	s.cache.Del(productID)
	return nil
}

func (s SchemaRepo) Delete(ctx context.Context, productID string) error {
	s.cache.Del(productID)
	err := s.db.DeleteWithFilter(ctx, mysql.ProductSchemaFilter{ProductID: productID})
	if err != nil {
		logx.WithContext(ctx).Errorf("%s.Delete err:%v", utils.FuncName(), err)
		return errors.Database
	}
	return err
}
