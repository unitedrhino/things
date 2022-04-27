package mysql

import (
	"context"
	"encoding/json"
	"github.com/dgraph-io/ristretto"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/domain/templateModel"
	"time"
)

const (
	expirtTime = time.Hour
)

type TemplateRepo struct {
	db    ProductTemplateModel
	cache *ristretto.Cache
}

func NewTemplateRepo(t ProductTemplateModel) templateModel.TemplateRepo {
	cache, _ := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	return &TemplateRepo{
		db:    t,
		cache: cache,
	}
}

func (t TemplateRepo) Insert(ctx context.Context, productID string, template *templateModel.Template) error {
	templateStr, err := json.Marshal(template)
	if err != nil {
		return errors.Parameter.WithMsg("模板的json格式不对")
	}
	_, err = t.db.Insert(&ProductTemplate{
		ProductID:   productID,
		Template:    string(templateStr),
		CreatedTime: time.Now(),
	})
	t.cache.SetWithTTL(productID, template, 1, expirtTime)
	return err
}

func (t TemplateRepo) GetTemplateInfo(ctx context.Context, productID string) (*templateModel.TemplateInfo, error) {
	temp, err := t.db.FindOne(productID)
	if err != nil {
		return nil, err
	}
	return &templateModel.TemplateInfo{
		ProductID:   temp.ProductID,
		Template:    temp.Template,
		CreatedTime: temp.CreatedTime,
	}, nil
}

func (t TemplateRepo) GetTemplate(ctx context.Context, productID string) (*templateModel.Template, error) {
	temp, ok := t.cache.Get(productID)
	if ok {
		return temp.(*templateModel.Template), nil
	}
	templateInfo, err := t.db.FindOne(productID)
	if err != nil {
		return nil, err
	}
	tempModel, err := templateModel.NewTemplate([]byte(templateInfo.Template))
	if err != nil {
		return nil, err
	}
	t.cache.SetWithTTL(productID, tempModel, 1, expirtTime)
	return tempModel, nil
}

func (t TemplateRepo) Update(ctx context.Context, productID string, template *templateModel.Template) error {
	templateStr, err := json.Marshal(template)
	if err != nil {
		return errors.Parameter.WithMsg("模板的json格式不对")
	}
	t.cache.Del(productID)
	err = t.db.Update(&ProductTemplate{
		ProductID:   productID,
		Template:    string(templateStr),
		CreatedTime: time.Now(),
	})
	return err
}

func (t TemplateRepo) ClearCache(ctx context.Context, productID string) error {
	t.cache.Del(productID)
	return nil
}

func (t TemplateRepo) Delete(ctx context.Context, productID string) error {
	t.cache.Del(productID)
	err := t.db.Delete(productID)
	return err
}
