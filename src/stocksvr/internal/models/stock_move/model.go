package stock_move

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/i-Things/things/shared/tools/ecode"
	"github.com/i-Things/things/shared/xzero/xmodel"
	"github.com/i-Things/things/shared/xzero/xsql"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type (
	CacheModel interface {
		xmodel.ICacheModel
		// CacheGetEntity 通过缓存拿单条数据
		CacheGetEntity(id uint64, checkExist bool) (err error, entity *CacheEntity)
		// QueryGetPkEntity 通过主键查询单条
		QueryGetPkEntity(id uint64, checkExist bool) (err error, entity *Entity)
		// QueryEntities 查询多条
		QueryEntities(selectBuilder squirrel.SelectBuilder) (err error, entity Entities)

		//QueryEntity 自定义查询单条数据 请自行组装条件
		QueryEntity(selectBuilder squirrel.SelectBuilder, checkExist bool) (err error, entity *Entity)
		FormatPrimary(primary interface{}) string
		Tx(session sqlx.Session) CacheModel
		Create(dto DtoCreate) (res int64, err error)
		UpdateDto(dto Dto) (err error)
	}
	defaultCacheModel struct {
		xmodel.CacheModel
		table string
	}
)

func NewCacheModel(conn sqlx.SqlConn, c cache.CacheConf) CacheModel {
	return &defaultCacheModel{
		CacheModel: xmodel.CacheModel{
			CachedConn: xsql.NewConn(conn, c),
		},
		table: Table,
	}
}

func (m *defaultCacheModel) Tx(session sqlx.Session) CacheModel {
	return &defaultCacheModel{
		CacheModel: xmodel.CacheModel{
			CachedConn: m.NewWithSession(session),
		},
		table: Table,
	}
}

// NewSelectBuilder 查询构造器
func (m *defaultCacheModel) NewSelectBuilder() squirrel.SelectBuilder {
	return squirrel.Select().From(m.table)
}

func (m *defaultCacheModel) CacheGetEntity(id uint64, checkExist bool) (error, *CacheEntity) {
	cacheIdKey := m.FormatPrimary(id)
	queryBuilder := m.queryPk(id)
	var entity CacheEntity
	err := m.CacheEntityWithPointer(cacheIdKey, queryBuilder, &entity, CacheEntityFieldNames)
	switch err {
	case nil:
		return nil, &entity
	case sql.ErrNoRows:
		if checkExist {
			return ecode.New(184764, "数据不存在"), nil
		} else {
			return nil, nil
		}
	default:
		return ecode.DB(19847, err), nil
	}
}

// 默认查询一个数据条件
func (m *defaultCacheModel) queryPk(id uint64) squirrel.SelectBuilder {
	return m.NewSelectBuilder().Where(" `id` = ? ", id).Limit(1)
}

// QueryGetPkEntity 通过主键查询单条数据
func (m *defaultCacheModel) QueryGetPkEntity(id uint64, checkExist bool) (error, *Entity) {
	queryBuilder := m.queryPk(id)
	return m.QueryEntity(queryBuilder, checkExist)
}

// QueryEntity 自己组装查询单条数据 如 QueryGetPkEntity
func (m *defaultCacheModel) QueryEntity(selectBuilder squirrel.SelectBuilder, checkExist bool) (error, *Entity) {
	var entity Entity
	err, exist := m.SelectEntityWithPointer(selectBuilder, &entity, EntityFieldNames)
	if err != nil {
		return ecode.DB(19887, err), nil
	}
	if !exist {
		if checkExist {
			return ecode.New(184766, "数据不存在"), nil
		} else {
			return nil, nil
		}
	}

	return nil, &entity
}

// QueryEntities 查询多条数据
func (m *defaultCacheModel) QueryEntities(queryBuilder squirrel.SelectBuilder) (err error, entity Entities) {
	err = m.SelectEntitiesWithPointer(queryBuilder, &entity, EntityFieldNames)
	return
}
func (m *defaultCacheModel) Insert(data map[string]interface{}, cacheKey ...string) (sql.Result, error) {
	//data["created_at"] = gtime.New().Format("Y-m-d H:i:s")
	return m.Exec(func(conn sqlx.Session) (result sql.Result, err error) {
		query, args, _ := squirrel.Insert(m.table).SetMap(data).ToSql()
		return conn.Exec(query, args...)
	}, cacheKey...)
}

func (m *defaultCacheModel) BatchInsert(data []map[string]interface{}, cacheKey ...string) (sql.Result, error) {
	insertBuilder := squirrel.Insert(m.table)
	insertBuilder = xmodel.InsertSetMaps(insertBuilder, data)
	return m.Exec(func(conn sqlx.Session) (result sql.Result, err error) {
		query, args, _ := insertBuilder.ToSql()
		return conn.Exec(query, args...)
	}, cacheKey...)
}

func (m *defaultCacheModel) NewInsertBuilder() squirrel.InsertBuilder {
	return squirrel.Insert(m.table)
}

// NewUpdateBuilder 更新构造器
func (m *defaultCacheModel) NewUpdateBuilder() squirrel.UpdateBuilder {
	return squirrel.Update(m.table)
}

// Update 通过更新构造器更新语句
// TODO: 更新时候注意 更新时间
func (m *defaultCacheModel) Update(updateBuilder squirrel.UpdateBuilder, data map[string]interface{}, cacheKey ...string) error {
	//data["updated_at"] = gtime.New().Format("Y-m-d H:i:s")
	_, err := m.Exec(func(conn sqlx.Session) (result sql.Result, err error) {
		query, args, _ := updateBuilder.Table(m.table).SetMap(data).ToSql()
		return conn.Exec(query, args...)
	}, cacheKey...)
	return err
}

// NewDeleteBuilder 删除构造器
func (m *defaultCacheModel) NewDeleteBuilder() squirrel.DeleteBuilder {
	return squirrel.Delete(m.table)
}

// Delete 通过构造器删除数据
func (m *defaultCacheModel) Delete(deleteBuilder squirrel.DeleteBuilder, cacheKey ...string) error {
	//data["deleted_at"] = gtime.New().Format("Y-m-d H:i:s")
	_, err := m.Exec(func(conn sqlx.Session) (result sql.Result, err error) {
		query, args, _ := deleteBuilder.ToSql()
		return conn.Exec(query, args...)
	}, cacheKey...)
	return err
}

func (m *defaultCacheModel) FormatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheEntityIdPrefix, primary)
}
