package relationDB

import "gitee.com/i-Things/share/stores"

// 示例
type ViewExample struct {
	ID int64 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"` // id编号
}

type ViewProjectInfo struct {
	ID int64 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"` // id编号
	//TenantCode    stores.TenantCode `gorm:"column:tenant_code;uniqueIndex:ri_mi;type:VARCHAR(50);NOT NULL"` // 租户编码
	//ProjectID     stores.ProjectID  `gorm:"column:projectID;type:bigint;NOT NULL"`                          // 所属项目ID(雪花ID)
	IndexImage    string `gorm:"column:index_image;type:varchar(200)"`        //图片地址
	Name          string `gorm:"column:name;type:varchar(50)"`                //项目名称
	Desc          string `gorm:"column:desc;type:varchar(200)"`               //项目描述
	CreatedUserID int64  `gorm:"column:created_user_id;type:bigint;NOT NULL"` //创建者id
	Status        int64  `gorm:"column:status;type:SMALLINT;default:1"`       //项目状态 1: 已发布 2: 未发布
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;index"`
	Detail      *ViewProjectDetail `gorm:"foreignKey:ProjectID;references:ID"`
}

type ViewProjectDetail struct {
	ID        int64  `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"` // id编号
	Content   string `gorm:"column:content;type:text"`                         //项目参数
	ProjectID int64  `gorm:"column:project_id;type:bigint;NOT NULL"`           //所属项目ID(雪花ID)
	stores.NoDelTime
	DeletedTime stores.DeletedTime `gorm:"column:deleted_time;index"`
}
