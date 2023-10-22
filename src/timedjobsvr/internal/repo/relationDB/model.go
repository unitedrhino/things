package relationDB

import (
	"github.com/i-Things/things/src/timedjobsvr/internal/domain"
	"time"
)

// 示例
type TimedExample struct {
	ID int64 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"` // id编号
}

type TimedJobLog struct {
	ID          int64               `gorm:"column:id;primary_key"`                     // 任务ID
	Group       string              `gorm:"column:group"`                              // 任务组名
	Type        string              `gorm:"column:type"`                               //任务类型:queue(消息队列消息发送)  sql(执行sql) email(邮件发送) http(http请求)
	SubType     string              `gorm:"column:sub_type;default:''"`                //任务子类型 natsJs nats
	Name        string              `gorm:"column:name"`                               // 任务名称
	Code        string              `gorm:"column:code"`                               //任务编码
	ResultCode  int64               `gorm:"column:return_code"`                        //结果code
	ResultMsg   string              `gorm:"column:return_msg"`                         //结果消息
	ExecLog     []*domain.ScriptLog `gorm:"column:exec_log;type:json;serializer:json"` //执行日志
	CreatedTime time.Time           `gorm:"column:created_time;index:,sort:desc;default:CURRENT_TIMESTAMP;NOT NULL"`
}

func (t *TimedJobLog) TableName() string {
	return "timed_job_log"
}
