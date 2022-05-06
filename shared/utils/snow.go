package utils

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

/*
 * twitter雪花算法golang实现,生成唯一趋势自增id
 * 保留位:63位
 * 毫秒时间戳:[62-20]43位,时间范围[1970-01-01 00:00:00.000,2248-09-26 15:10:22.207]
 * 机器id:[19-12]8位,十进制范围[0,255]
 * 序列号:[11-0]12位,十进制范围[0,4095]
 * bobo
 */
const epoch = int64(1609430400) // 设置起始时间(时间戳/毫秒)：2020-01-01 00:00:00，有效期69年
type SnowFlake struct {
	machineID int64      //机器 id占8位,十进制范围是[0,255]
	sn        int64      //序列号占12位,十进制范围是[0,4095]
	lastTime  int64      //上次的时间戳(毫秒级)
	_lock     sync.Mutex //锁
}

func GetNodeID(cache cache.ClusterConf, svrName string) int64 {
	key := fmt.Sprintf("node:id:%s", svrName)
	nodeIdS, err := kv.NewStore(cache).Incr(key)
	if err != nil {
		nodeIdS = rand.NewSource(time.Now().UnixNano()).Int63()
	}
	return nodeIdS % 1024
}

//var Snow = &SnowFlake{
//	lastTime: time.Now().UnixNano() / 1000000,
//	machineID:1,
//}
func NewSnowFlake(mId int64) *SnowFlake {
	sf := SnowFlake{
		lastTime: time.Now().UnixNano() / 1000000,
	}
	sf.SetMachineId(mId)
	return &sf
}

func (c *SnowFlake) lock() {
	c._lock.Lock()
}

func (c *SnowFlake) unLock() {
	c._lock.Unlock()
}

//获取当前毫秒
func (c *SnowFlake) getCurMilliSecond() int64 {
	return time.Now().UnixNano() / 1000000
}

//设置机器id,默认为0,范围[0,255]
func (c *SnowFlake) SetMachineId(mId int64) {
	//保留8位
	mId = mId & 0xFF
	//左移12位,序列号是12位的
	mId <<= 12
	c.machineID = mId
}

//获取机器id
func (c *SnowFlake) GetMachineId() int64 {
	mId := c.machineID
	mId >>= 12
	return mId | 0xFF
}

//解析雪花(id)
// 返回值
// milliSecond:毫秒数
// mId:机器id
// sn:序列号
func (c *SnowFlake) ParseId(id int64) (milliSecond, mId, sn int64) {
	sn = id & 0xFFF
	id >>= 12
	mId = id & 0xFF
	id >>= 8
	milliSecond = id & 0x7FFFFFFFFFF

	return
}

//毫秒转换成time
func (c *SnowFlake) MilliSecondToTime(milliSecond int64) (t time.Time) {
	return time.Unix(milliSecond/1000, milliSecond%1000*1000000)
}

//毫秒转换成"20060102T150405.999Z"
func (c *SnowFlake) MillisecondToTimeTz(ts int64) string {
	tm := c.MilliSecondToTime(ts)
	return tm.UTC().Format("20060102T150405.999Z")
}

//毫秒转换成"2006-01-02 15:04:05.999"
func (c *SnowFlake) MillisecondToTimeDb(ts int64) string {
	tm := c.MilliSecondToTime(ts)
	return tm.UTC().Format("2006-01-02 15:04:05.999")
}

//获取雪花
//返回值
//id:自增id
//ts:生成该id的毫秒时间戳
func (c *SnowFlake) GetSnowflakeId() (id int64) {
	curTime := c.getCurMilliSecond()
	var sn int64 = 0

	c.lock()
	// 同一毫秒
	if curTime == c.lastTime {
		c.sn++
		// 序列号占 12 位,十进制范围是 [0,4095]
		if c.sn > 4095 {
			for {
				// 让出当前线程
				runtime.Gosched()
				curTime = c.getCurMilliSecond()
				if curTime != c.lastTime {
					break
				}
			}
			c.sn = 0
		}
	} else {
		c.sn = 0
	}
	sn = c.sn
	c.lastTime = curTime
	c.unLock()

	//当前时间小于上次的时间，系统时间改过了吗?
	/*
	   if curTimeStamp < c.lastTimeStamp {
	           return 0, curTimeStamp
	   }
	*/
	//机器id占用8位空间,序列号占用12位空间,所以左移20位
	rightBinValue := (curTime - epoch) & 0x7FFFFFFFFFF
	rightBinValue <<= 20
	id = rightBinValue | c.machineID | sn

	return id
}
