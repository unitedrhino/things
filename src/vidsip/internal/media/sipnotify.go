package media

import (
	"fmt"
	"gitee.com/i-Things/core/shared/utils"
	db "github.com/i-Things/things/src/vidsip/internal/repo/relationDB"
	"github.com/sirupsen/logrus"
	"net/url"
	"strings"
	"time"
)

func notify(data *Notify) {
	if url, ok := NotifyMap[data.Method]; ok {
		res, err := utils.PostJSONRequest(url, data)
		if err != nil {
			logrus.Warningln(data.Method, "send notify fail.", err)
		}
		if strings.ToUpper(string(res)) != "OK" {
			logrus.Warningln(data.Method, "send notify resp fail.", string(res), "len:", len(res), NotifyMap, data)
		} else {
			logrus.Debug("notify send succ:", data.Method, data.Data)
		}
	} else {
		logrus.Traceln("notify config not found", data.Method)
	}
}

func notifyDevicesAcitve(id, status string) *Notify {
	return &Notify{
		Method: NotifyMethodDevicesActive,
		Data: map[string]interface{}{
			"deviceid": id,
			"status":   status,
			"time":     time.Now().Unix(),
		},
	}
}
func notifyDevicesRegister(u db.SipDevices) *Notify {
	u.Sys = SipInfo
	return &Notify{
		Method: NotifyMethodDevicesRegister,
		Data:   u,
	}
}

func notifyChannelsActive(d db.SipChannels) *Notify {
	return &Notify{
		Method: NotifyMethodChannelsActive,
		Data: map[string]interface{}{
			"channelid": d.ChannelID,
			"status":    d.Status,
			"time":      time.Now().Unix(),
		},
	}
}
func notifyRecordStop(url string, req url.Values) *Notify {
	d := map[string]interface{}{
		"url": fmt.Sprintf("%s/%s", "http://192.168.10.117:8088", url),
	}
	for k, v := range req {
		d[k] = v[0]
	}
	return &Notify{
		Method: NotifyMethodRecordStop,
		Data:   d,
	}
}
