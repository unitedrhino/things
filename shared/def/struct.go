package def

type PageInfo struct {
	Page       int64  `json:"page" form:"page"`             // 页码
	PageSize   int64  `json:"pageSize" form:"pageSize"`     // 每页大小
	SearchKey  string `json:"searchKey" form:"searchKey"`   // 搜索的key
	SearchType string `json:"searchType" form:"searchType"` // 搜索的类型
}

func (p PageInfo) GetLimit() int64 {
	if p.PageSize == 0 {
		return 20
	}
	return p.PageSize
}
func (p PageInfo) GetOffset() int64 {
	if p.Page == 0 {
		return 0
	}
	return p.PageSize * (p.Page - 1)
}
