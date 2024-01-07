package def

import (
	"fmt"
	"gorm.io/gorm"
	"time"

	sq "github.com/Masterminds/squirrel"
)

const (
	OrderAes  = iota //从久到近排序
	OrderDesc        //时间从近到久排序
)

var orderMap = map[int64]string{
	OrderAes:  "aes",
	OrderDesc: "desc",
}

type ExecArgs struct {
	Query string
	Args  []any
}

type PageInfo struct {
	Page   int64     `json:"page" form:"page"`         // 页码
	Size   int64     `json:"pageSize" form:"pageSize"` // 每页大小
	Orders []OrderBy `json:"orderBy" form:"orderBy"`   // 排序信息
}
type PageInfo2 struct {
	TimeStart int64     `json:"timeStart"`
	TimeEnd   int64     `json:"timeEnd"`
	Page      int64     `json:"page" form:"page"`       // 页码
	Size      int64     `json:"size" form:"size"`       // 每页大小
	Orders    []OrderBy `json:"orderBy" form:"orderBy"` // 排序信息
}

// 排序结构体
type OrderBy struct {
	Filed string `json:"filed" form:"filed"` //要排序的字段名
	Sort  int64  `json:"sort" form:"sort"`   //排序的方式：0 OrderAes、1 OrderDesc
}

type TimeRange struct {
	Start int64
	End   int64
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

// 获取排序参数
func (p *PageInfo) GetOrders() (arr []string) {
	if p != nil && len(p.Orders) > 0 {
		for _, o := range p.Orders {
			arr = append(arr, fmt.Sprintf("%s %s", o.Filed, orderMap[o.Sort]))
		}
	}
	return
}
func (p *PageInfo) ToGorm(db *gorm.DB) *gorm.DB {
	if p == nil {
		return db
	}
	db = db.Offset(int(p.GetOffset())).Limit(int(p.GetLimit()))
	if len(p.Orders) != 0 {
		orders := p.GetOrders()
		for _, o := range orders {
			db = db.Order(o)
		}
	} else {
		db.Order("created_time desc")
	}
	return db
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

func (t TimeRange) FmtSql(sql sq.SelectBuilder) sq.SelectBuilder {
	if t.Start != 0 {
		sql = sql.Where("created_time>=?", time.Unix(t.Start, 0))
	}
	if t.End != 0 {
		sql = sql.Where("created_time<=?", time.Unix(t.End, 0))
	}
	return sql
}
