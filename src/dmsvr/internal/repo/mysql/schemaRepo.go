package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/dgraph-io/ristretto"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/domain/schema"
	"time"
)

const (
	expirtTime = time.Hour
)

type SchemaRepo struct {
	db    productSchemaModel
	cache *ristretto.Cache
}

func NewSchemaRepo(t productSchemaModel) schema.SchemaRepo {
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

func (t SchemaRepo) Insert(ctx context.Context, productID string, schema *schema.Model) error {
	templateStr, err := json.Marshal(schema)
	if err != nil {
		return errors.Parameter.WithMsg("模板的json格式不对")
	}
	_, err = t.db.Insert(ctx, &ProductSchema{
		ProductID:   productID,
		Schema:      string(templateStr),
		CreatedTime: time.Now(),
	})
	t.cache.SetWithTTL(productID, schema, 1, expirtTime)
	return err
}

func (t SchemaRepo) GetSchemaInfo(ctx context.Context, productID string) (*schema.SchemaInfo, error) {
	temp, err := t.db.FindOne(ctx, productID)
	if err != nil {
		return nil, err
	}
	return &schema.SchemaInfo{
		ProductID:   temp.ProductID,
		Template:    temp.Schema,
		CreatedTime: temp.CreatedTime,
	}, nil
}

func (t SchemaRepo) GetSchemaModel(ctx context.Context, productID string) (*schema.Model, error) {
	temp, ok := t.cache.Get(productID)
	if ok {
		return temp.(*schema.Model), nil
	}
	templateInfo, err := t.db.FindOne(ctx, productID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Parameter.AddMsg("ProductID not find")
		}
		return nil, err
	}
	tempModel, err := schema.NewSchema([]byte(templateInfo.Schema))
	if err != nil {
		return nil, err
	}
	t.cache.SetWithTTL(productID, tempModel, 1, expirtTime)
	return tempModel, nil
}

func (t SchemaRepo) Update(ctx context.Context, productID string, schema *schema.Model) error {
	templateStr, err := json.Marshal(schema)
	if err != nil {
		return errors.Parameter.WithMsg("模板的json格式不对")
	}
	t.cache.Del(productID)
	old, err := t.db.FindOne(ctx, productID)
	if err != nil {
		return errors.Database
	}
	err = t.db.Update(ctx, &ProductSchema{
		ProductID:   productID,
		Schema:      string(templateStr),
		UpdatedTime: time.Now(),
		CreatedTime: old.CreatedTime,
	})
	return err
}

func (t SchemaRepo) ClearCache(ctx context.Context, productID string) error {
	t.cache.Del(productID)
	return nil
}

func (t SchemaRepo) Delete(ctx context.Context, productID string) error {
	t.cache.Del(productID)
	err := t.db.Delete(ctx, productID)
	return err
}
