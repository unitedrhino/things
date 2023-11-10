package relationDB

import (
	"database/sql"
	"github.com/i-Things/things/shared/stores"
)

type DmExample struct {
	ID int64 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"` // id编号
}

// 服务信息
type VidmgrInfo struct {
	VidmgrID     string            `gorm:"column:id;type:char(11);primary_key;NOT NULL"`                // 服务id
	VidmgrName   string            `gorm:"column:name;type:varchar(100);NOT NULL"`                      // 服务名称
	VidmgrIpV4   int64             `gorm:"column:ipv4;type:bigint"`                                     // 服务IP
	VidmgrPort   int64             `gorm:"column:port;type:bigint"`                                     // 服务端口
	VidmgrType   int64             `gorm:"column:type;type:smallint;default:1"`                         // 服务类型:1:zlmediakit,2:srs,3:monibuca
	VidmgrStatus int64             `gorm:"column:status;type:smallint;default:0;NOT NULL"`              //服务状态: 0：未激活 1：在线  2:离线
	VidmgrSecret string            `gorm:"column:secret;type:varchar(50)"`                              // 服务秘钥
	FirstLogin   sql.NullTime      `gorm:"column:first_login"`                                          // 激活时间
	LastLogin    sql.NullTime      `gorm:"column:last_login"`                                           // 最后登录时间
	Desc         string            `gorm:"column:desc;type:varchar(200)"`                               // 描述
	Tags         map[string]string `gorm:"column:tags;type:json;serializer:json;NOT NULL;default:'{}'"` // 产品标签
	stores.Time
}

func (m *VidmgrInfo) TableName() string {
	return "vid_mgr_info"
}

// 视频流信息
type VidstreamInfo struct {
	StreamID       string            `gorm:"column:id;type:bigint;primary_key;NOT NULL"` // 视频流的id
	StreamName     string            `gorm:"column:name;type:varchar(100);NOT NULL"`     // 服务名称
	NetType        int64             `gorm:"column:net_type;type:smallint;NOT NULL"`     // 服务名称
	DevType        int64             `gorm:"column:dev_type;type:smallint;default:1"`
	DevStreamType  int64             `gorm:"column:dev_streamtype;type:smallint;default:1"`
	ChannelID      string            `gorm:"column:channel_id;type:varchar(32)"`
	ChannelName    string            `gorm:"column:channel_name;type:varchar(32)"`
	LowNetType     int64             `gorm:"column:low_nettype;type:smallint;default:1"`
	IsShareChannel bool              `gorm:"column:share_channel;type:bit(1);default:0;NOT NULL"`
	IsAutoPush     bool              `gorm:"column:auto_push;type:bit(1);default:0;NOT NULL"`
	IsAutoRecord   bool              `gorm:"column:auto_record;type:bit(1);default:0;NOT NULL"`
	IsPTZ          bool              `gorm:"column:is_ptz;type:bit(1);default:0;NOT NULL"`
	IsOnline       bool              `gorm:"column:is_online;type:bit(1);default:0;NOT NULL"`
	VidmgrInfo     *VidmgrInfo       `gorm:"foreignKey:VidmgrID;references:VidmgrID"`                     // 添加外键
	Desc           string            `gorm:"column:desc;type:varchar(200)"`                               // 描述
	Tags           map[string]string `gorm:"column:tags;type:json;serializer:json;NOT NULL;default:'{}'"` // 产品标签
	stores.Time
}

func (m *VidstreamInfo) TableName() string {
	return "vid_stream_info"
}
