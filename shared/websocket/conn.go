package websocket

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"gitee.com/asktop_golib/util/aslice"
	"github.com/gorilla/websocket"
	"github.com/hashicorp/go-uuid"
	e "github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

var (
	dp *dispatcher //ws调度器
)

const (
	errorCount    = 5                     //错误次数
	interval      = 5 * time.Second       //心跳间隔
	keepAliveType = websocket.PingMessage //心跳类型
)

type connection struct {
	r        *http.Request
	server   *Server
	ws       *websocket.Conn   //ws连接实例
	clientId string            //ws连接实例唯一标识
	closed   bool              //ws连接已关闭
	send     chan []byte       //发送信息管道
	topics   map[string]string //订阅信息
	pingErrs []int64           //发送的心跳失败次数
	pongErrs []int64           //收到的心跳失败次数
}

// ws调度器
type dispatcher struct {
	s2cGzip  bool                              //发送的信息是否gzip压缩
	connPool map[string]*connection            //ws连接池 map[clientId]*connection
	subPool  map[string]map[string]*connection //订阅池 map[topic]map[clientId]*connection
	sendSub  chan WsResp                       //发送订阅
	mu       sync.Mutex                        // 互斥锁
}

// 创建ws调度器
func StartWsDp(s2cGzip bool) {
	if dp == nil {
		dp = newDp(s2cGzip)
	}
}

// 创建ws调度器
func newDp(s2cGzip bool) *dispatcher {
	d := &dispatcher{
		s2cGzip:  s2cGzip,
		connPool: make(map[string]*connection),
		subPool:  make(map[string]map[string]*connection),
		sendSub:  make(chan WsResp, 10000),
	}
	go func(d *dispatcher) {
		for {
			select {
			//发送订阅
			case subBody := <-d.sendSub:
				subs, ok := d.subPool[subBody.Path]
				if !ok {
					break
				}
				d.mu.Lock()
				for _, conn := range subs {
					conn.sendMessage(subBody)
				}
				d.mu.Unlock()
			}
		}
	}(d)
	return d
}

// 读ping心跳
func (c *connection) pingRead(message []byte) error {
	logx.Infof("%s.[ws] message:%s clientId:%s", utils.FuncName(), string(message), c.clientId)
	if aslice.ContainInt64(c.pingErrs, int64(binary.BigEndian.Uint64(message))) {
		c.pingErrs = []int64{}
	} else {
		c.writeMessage(websocket.PingMessage, []byte("ping error message :"+string(message)))
	}
	return nil
}

// 读pong心跳
func (c *connection) pongRead(message []byte) error {
	logx.Infof("%s.[ws] message:%s clientId:%s", utils.FuncName(), string(message), c.clientId)
	if aslice.ContainInt64(c.pongErrs, int64(binary.BigEndian.Uint64(message))) {
		c.pongErrs = []int64{}
	} else {
		c.writeMessage(websocket.PongMessage, []byte("pong error message :"+string(message)))
	}
	return nil
}

// 发送ping心跳
func (c *connection) pingSend() error {
	if len(c.pingErrs) >= errorCount || len(c.pongErrs) >= errorCount {
		//连续5次没有收到ping心跳 关闭连接
		return errors.New("connection timeout")
	}
	nowTime := []byte(strconv.FormatInt(time.Now().Unix(), 10))
	if err := c.writeMessage(websocket.PingMessage, nowTime); err != nil {
		c.pingErrs = append(c.pingErrs, int64(binary.BigEndian.Uint64(nowTime)))
	} else {
		c.pingErrs = []int64{}
		c.pongErrs = append(c.pongErrs, int64(binary.BigEndian.Uint64(nowTime)))
	}
	return nil
}

// 发送pong心跳
func (c *connection) pongSend() error {
	if len(c.pingErrs) >= errorCount || len(c.pongErrs) >= errorCount {
		//连续5次没有收到pong心跳 关闭连接
		return errors.New("connection timeout")
	}
	nowTime := []byte(strconv.FormatInt(time.Now().Unix(), 10))
	if err := c.writeMessage(websocket.PongMessage, nowTime); err != nil {
		c.pongErrs = append(c.pongErrs, int64(binary.BigEndian.Uint64(nowTime)))
	} else {
		c.pongErrs = []int64{}
		c.pingErrs = append(c.pingErrs, int64(binary.BigEndian.Uint64(nowTime)))
	}
	return nil
}

// 发送订阅信息
func SendSub(body WsResp) {
	//判断
	dp.sendSub <- body
}

// 创建ws连接
func NewConn(server *Server, r *http.Request, wsConn *websocket.Conn) *connection {
	var clientId string
	for {
		clientId, _ = uuid.GenerateUUID()
		if _, ok := dp.connPool[clientId]; !ok {
			break
		}
	}
	conn := &connection{
		server:   server,
		ws:       wsConn,
		r:        r,
		clientId: clientId,
		send:     make(chan []byte, 10000),
		topics:   make(map[string]string),
	}
	dp.connPool[clientId] = conn
	logx.Infof("%s.[ws]创建连接成功 RemoteAddr::%s clientId:%s", utils.FuncName(), wsConn.RemoteAddr().String(), clientId)
	resp := WsResp{StatusCode: http.StatusOK}
	conn.sendMessage(resp)
	return conn
}

// 开启读取进程
func (c *connection) StartRead() {
	defer func() {
		c.Close("read message error")
	}()
	c.ws.SetPongHandler(func(message string) error {
		c.pongRead([]byte(message))
		return nil
	})
	c.ws.SetPingHandler(func(message string) error {
		c.pingRead([]byte(message))
		return nil
	})
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		logx.Infof("%s.[ws] message:%s clientId:%s", utils.FuncName(), string(message), c.clientId)
		var data map[string]interface{}
		err = json.Unmarshal(message, &data)
		if err != nil {
			c.errorSend(e.Type.AddDetail("error reading message"))
			continue
		}
		c.handleRequest(message)
	}
}
func (c *connection) errorSend(data error) {
	resp := WsResp{
		StatusCode: http.StatusBadRequest,
		WsBody:     WsBody{Body: data},
	}
	c.sendMessage(resp)
}

func (c *connection) handleRequest(message []byte) {
	var body WsReq
	err := json.Unmarshal(message, &body)
	if err != nil {
		c.errorSend(e.Parameter)
		return
	}
	if err := isDataComplete(body.Type, body); err != nil {
		c.errorSend(err)
		return
	}
	switch body.Type {
	case Control:
		downControl(c, body)
	case Sub:
		subscribeHandle(c, body)
	case UnSub:
		unSubscribeHandle(c, body)
	default:
	}
}

func isDataComplete(wsType WsType, body WsReq) error {
	if wsType == "" {
		return e.Parameter.AddDetail("type is  null")
	}
	switch wsType {
	case Control:
		if body.Path == "" || body.Method == "" || body.Body == "" {
			return e.Parameter.AddDetail("path|method|body is  null")
		}
	case Sub, UnSub:
		if body.Path == "" {
			return e.Parameter.AddDetail("path  is  null")
		}
	case Pub:
		return e.NotRealize
	default:
	}
	return nil
}

func downControl(c *connection, body WsReq) {
	reqBody, err := getRequestBody(body.Body)
	if err != nil {
		// 处理编码错误
	}
	bodyBytes, err := json.Marshal(body.Body)
	length := len(bodyBytes)
	header := c.r.Header
	header.Set("Content-Type", "application/json")
	header.Set("Content-Length", strconv.Itoa(length))
	r := &http.Request{
		Method: body.Method,
		Host:   c.r.Host,
		URL: &url.URL{
			Path: body.Path,
		},
		Header:        header,
		Body:          reqBody,
		ContentLength: int64(length),
	}
	w := response{req: &body, resp: WsResp{WsBody: WsBody{Handler: map[string][]string{}, Type: ControlRet}}}
	c.server.ServeHTTP(&w, r)
	c.sendMessage(w.resp)
}

// 将请求体转换为io.ReadCloser类型
func getRequestBody(body interface{}) (io.ReadCloser, error) {
	var reqBody io.ReadCloser
	if body != nil {
		switch body.(type) {
		case string:
			reqBody = ioutil.NopCloser(bytes.NewBufferString(body.(string)))
		case []byte:
			reqBody = ioutil.NopCloser(bytes.NewBuffer(body.([]byte)))
		case map[string]interface{}:
			bodyBytes, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			reqBody = ioutil.NopCloser(bytes.NewReader(bodyBytes))
		default:
			// 处理其他类型
		}
	}
	return reqBody, nil
}

// 开启发送进程
func (c *connection) StartWrite() {
	ticker := time.NewTicker(interval)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		//发送心跳
		case <-ticker.C:
			if c.closed {
				return
			}
			var err error
			switch keepAliveType {
			case websocket.PingMessage:
				err = c.pingSend()
			case websocket.PongMessage:
				err = c.pongSend()
			}
			if err != nil {
				c.Close("connection timeout")
				return
			}
		//发送信息
		case message := <-c.send:
			if c.closed {
				return
			}
			if err := c.writeMessage(websocket.TextMessage, message); err != nil {
				c.Close("send message error")
				return
			}
		}
	}
}

// 关闭ws连接
func (c *connection) Close(msg string) {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	_, ok := dp.connPool[c.clientId]
	if ok || !c.closed {
		c.closed = true
		close(c.send)
		delete(dp.connPool, c.clientId)
		for _, subs := range dp.subPool {
			delete(subs, c.clientId)
		}
		c.ws.Close()
		logx.Infof("%s.[ws]关闭连接  clientId:%s", utils.FuncName(), c.clientId)
	}
}

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

// 发送信息
func (c *connection) sendMessage(body WsResp) {
	message, _ := json.Marshal(body)
	if !c.closed {
		c.send <- message
	}
}

// 写消息
func (c *connection) writeMessage(messageType int, message []byte) error {
	if message == nil {
		logx.Infof("%s.[ws]error message: is  null ")
	}
	switch messageType {
	case websocket.PingMessage, websocket.PongMessage:
		err := c.ws.WriteControl(messageType, message, time.Time{})
		if err != nil {
			logx.Infof("%s.[ws]error message::%s clientId:%s", utils.FuncName(), string(message), c.clientId)
		}
	case websocket.TextMessage:
		err := c.ws.WriteMessage(messageType, message)
		if err != nil {
			logx.Infof("%s.[ws]error message::%s clientId:%s", utils.FuncName(), string(message), c.clientId)
		}
	}
	logx.Infof("%s.[ws] message:%s clientId:%s", utils.FuncName(), string(message), c.clientId)
	return nil
}
