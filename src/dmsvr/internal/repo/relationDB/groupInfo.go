package relationDB

import (
	"context"
	"database/sql"
	"github.com/i-Things/things/shared/store"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"gorm.io/gorm"
)

type GroupInfoRepo struct {
	db *gorm.DB
}

func NewGroupInfoRepo(in any) *GroupInfoRepo {
	return &GroupInfoRepo{db: store.GetCommonConn(in)}
}

func (g GroupInfoRepo) Insert(ctx context.Context, data *mysql.DmGroupInfo) (sql.Result, error) {
	//TODO implement me
	panic("implement me")
}

func (g GroupInfoRepo) FindOne(ctx context.Context, groupID int64) (*mysql.DmGroupInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (g GroupInfoRepo) FindOneByGroupName(ctx context.Context, groupName string) (*mysql.DmGroupInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (g GroupInfoRepo) Update(ctx context.Context, data *mysql.DmGroupInfo) error {
	//TODO implement me
	panic("implement me")
}

func (g GroupInfoRepo) Delete(ctx context.Context, groupID int64) error {
	//TODO implement me
	panic("implement me")
}
