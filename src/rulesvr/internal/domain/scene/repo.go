package scene

import "github.com/i-Things/things/shared/def"
import "context"

type InfoFilter struct {
	Name string `json:"name"`
}

type Repo interface {
	Insert(ctx context.Context, info *Info) error
	Update(ctx context.Context, info *Info) error
	Delete(ctx context.Context, id int64) error
	FindOne(ctx context.Context, id int64) (*Info, error)
	FindOneByName(ctx context.Context, name string) (*Info, error)
	FindByFilter(ctx context.Context, filter InfoFilter, page *def.PageInfo) ([]*Info, error)
	CountByFilter(ctx context.Context, filter InfoFilter) (size int64, err error)
}
