package relationDB

import "github.com/i-Things/things/shared/stores"

// 示例
type TimedQueueExample struct {
	ID int64 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"` // id编号
}

const (
	TaskStatusRun   = 1 //正常运行
	TaskStatusPause = 2 //暂停
	TaskStatusStop  = 3 //停用
)

type TimedTask struct {
	ID             int64  `gorm:"column:id;primary_key"`                                              // 任务ID
	Group          string `gorm:"column:group;uniqueIndex:uni_group_name;uniqueIndex:uni_group_code"` // 任务组名
	Type           string `gorm:"column:type"`                                                        //任务类型:queue(消息队列消息发送)  sql(执行sql) email(邮件发送) http(http请求)
	SubType        string `gorm:"column:sub_type;default:''"`                                         //任务子类型 natsJs nats
	Name           string `gorm:"column:name;uniqueIndex:uni_group_name"`                             // 任务名称
	Code           string `gorm:"column:code;uniqueIndex:uni_group_code"`                             //任务编码
	Params         string `gorm:"column:params;type:json;NOT NULL;default:'{}'"`                      // 任务参数
	CronExpression string `gorm:"column:cron_expression"`                                             // cron执行表达式
	Status         int64  `gorm:"column:status"`                                                      // 状态（1正常运行 2:暂停 3:停止使用）
	EntryID        string `gorm:"column:entry_id"`                                                    //执行任务的id
	Priority       string `gorm:"column:priority"`                                                    //优先级: 6:critical 最高优先级  3: default 普通优先级 1:low 低优先级
	stores.Time
}

func (t *TimedTask) TableName() string {
	return "timed_task"
}
