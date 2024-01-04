package websocket

import (
	e "github.com/i-Things/things/shared/errors"
	"net/http"
)

func subscribeHandle(c *connection, bd WsReq) {
	//校验订阅topic
	topic, err := wsCheckSub(bd.Path)
	var resp WsResp
	resp.WsBody = bd.WsBody
	resp.WsBody.Type = SubRet
	if err != nil {
		c.errorSend(e.Default.AddDetail("error  subscribe  topic"))
	} else {
		c.addSubscribe(topic) //添加订阅
		resp.StatusCode = http.StatusOK
	}
	c.sendMessage(resp)
}

// 处理取消订阅
func unSubscribeHandle(c *connection, bd WsReq) {
	//校验订阅topic
	topic, err := wsCheckSub(bd.Path)
	var resp WsResp
	resp.WsBody = bd.WsBody
	resp.WsBody.Type = UnSubRet
	if err != nil {
		c.errorSend(e.Default.AddDetail("error  unSubscribe  topic"))
		return
	} else {
		c.unSubscribe(topic) //取消订阅
		resp.StatusCode = http.StatusOK
		c.sendMessage(resp)
		return
	}
}

// 校验订阅topic合法性
func wsCheckSub(sub string) (topic string, err error) {
	//代码逻辑
	return sub, nil
}

// 添加订阅
func (c *connection) addSubscribe(topic string) {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	_, ok := dp.connPool[c.clientId]
	if !ok {
		return
	}
	subs, ok := dp.subPool[topic]
	if !ok {
		subs = make(map[string]*connection)
		dp.subPool[topic] = subs
	}
	subs[c.clientId] = c
	c.topics[topic] = topic
}

// 取消订阅
func (c *connection) unSubscribe(topic string) {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	delete(c.topics, topic)
	_, ok := dp.connPool[c.clientId]
	if !ok {
		return
	}
	subs, ok := dp.subPool[topic]
	if ok {
		delete(subs, c.clientId)
		if len(subs) == 0 {
			delete(dp.subPool, topic)
		}
	}
}
