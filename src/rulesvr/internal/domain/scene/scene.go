package scene

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"time"
)

type Info struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Desc        string    `json:"desc"`
	CreatedTime time.Time `json:"createdTime"`
	Trigger     *Trigger  `json:"trigger"`
	When        []*Term   `json:"when"`
	Then        []*Action `json:"then"`
}

type CreateInfoDto struct {
	Name        string    `json:"name"`
	Desc        string    `json:"desc"`
	createdTime time.Time `json:"createdTime"`
	Trigger     Trigger   `json:"trigger"`
	When        Term      `json:"when"`
	Action      Action    `json:"action"`
}

type InfoFilter struct {
}

type Repo interface {
	Insert(ctx context.Context, info *Info) error
	Update(ctx context.Context, info *Info) error
	Delete(ctx context.Context, id int64) error
	FindOne(ctx context.Context, id int64) (*Info, error)
	FindByFilter(ctx context.Context, filter InfoFilter, page *def.PageInfo) ([]*Info, error)
	CountByFilter(ctx context.Context, filter InfoFilter) (size int64, err error)
}
