package relationDB

import (
	"github.com/i-Things/things/shared/stores"
	"github.com/zeromicro/go-zero/core/logx"
)

func Migrate() error {
	db := stores.GetCommonConn(nil)
	if !db.Migrator().HasTable(&SysUserInfo{}) {
		//需要初始化表
		logx.Info("开始初始化sysvr的表")
	}
	return db.AutoMigrate(
		&SysUserInfo{},
		&SysRoleInfo{},
		&SysRoleMenu{},
		&SysMenuInfo{},
		&SysLoginLog{},
		&SysOperLog{},
		&SysApiInfo{},
		&SysApiAuth{},
	)
}
