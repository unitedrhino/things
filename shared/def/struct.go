package def

type PageInfo struct {
	Page     int64 `json:"page" form:"page"`         // 页码
	PageSize int64 `json:"pageSize" form:"pageSize"` // 每页大小
}

func (p PageInfo) GetLimit() int64 {
	return p.PageSize
}
func (p PageInfo) GetOffset() int64 {
	return p.PageSize * (p.Page - 1)
}
