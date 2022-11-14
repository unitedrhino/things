package def

import (
	"time"

	sq "github.com/Masterminds/squirrel"
)

const (
	OrderAes  = iota //从久到近排序
	OrderDesc        //时间从近到久排序
)

type PageInfo struct {
	Page int64 `json:"page" form:"page"`         // 页码
	Size int64 `json:"pageSize" form:"pageSize"` // 每页大小
}
type PageInfo2 struct {
	TimeStart int64 `json:"timeStart"`
	TimeEnd   int64 `json:"timeEnd"`
	Page      int64 `json:"page" form:"page"` // 页码
	Size      int64 `json:"size" form:"size"` // 每页大小
}

func (p *PageInfo) GetLimit() int64 {
	if p == nil || p.Size == 0 {
		return 2000
	}
	return p.Size
}
func (p *PageInfo) GetOffset() int64 {
	if p == nil || p.Page == 0 {
		return 0
	}
	return p.Size * (p.Page - 1)
}

func (p PageInfo2) GetLimit() int64 {
	//if p.Size == 0 {
	//	return 20000
	//}
	return p.Size
}
func (p PageInfo2) GetOffset() int64 {
	if p.Page == 0 {
		return 0
	}
	return p.Size * (p.Page - 1)
}
func (p PageInfo2) GetTimeStart() time.Time {
	return time.UnixMilli(p.TimeStart)
}
func (p PageInfo2) GetTimeEnd() time.Time {
	return time.UnixMilli(p.TimeEnd)
}

func (p PageInfo2) FmtSql(sql sq.SelectBuilder) sq.SelectBuilder {
	if p.TimeStart != 0 {
		sql = sql.Where("ts>=?", p.GetTimeStart())
	}
	if p.TimeEnd != 0 {
		sql = sql.Where("ts<=?", p.GetTimeEnd())
	}
	if p.Size != 0 {
		sql = sql.Limit(uint64(p.GetLimit()))
		if p.Page != 0 {
			sql = sql.Offset(uint64(p.GetOffset()))
		}
	}
	return sql
}

func (p PageInfo2) FmtWhere(sql sq.SelectBuilder) sq.SelectBuilder {
	if p.TimeStart != 0 {
		sql = sql.Where(sq.GtOrEq{"ts": p.GetTimeStart()})
	}
	if p.TimeEnd != 0 {
		sql = sql.Where(sq.LtOrEq{"ts": p.GetTimeEnd()})
	}
	return sql
}
