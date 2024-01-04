package sip

import (
	"bufio"
	"bytes"
	"github.com/i-Things/things/shared/utils"
	"io"
	"net"
	"strings"
	"time"
)

// Packet Packet
type Packet struct {
	reader     *bufio.Reader
	raddr      net.Addr
	bodylength int
}

func newPacket(data []byte, raddr net.Addr) Packet {
	//logrus.Traceln("receive new packet,from:", raddr.String(), string(data))
	return Packet{
		reader:     bufio.NewReader(bytes.NewReader(data)),
		raddr:      raddr,
		bodylength: getBodyLength(data),
	}
}

func (p *Packet) nextLine() (string, error) {
	str, err := p.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	if len(str) >= 2 {
		str = str[:len(str)-2]
	}
	return str, err
}

func (p *Packet) bodyLength() int {
	return p.bodylength
}

func (p *Packet) getBody() ([]byte, error) {
	if p.bodyLength() < 1 {
		return []byte{}, nil
	}
	body := make([]byte, p.bodylength)
	if p.bodylength > 0 {
		n, err := io.ReadFull(p.reader, body)
		if err != nil && err != io.ErrUnexpectedEOF {
			return body, err
		}
		if n != p.bodylength {
			//logrus.Warningf("body length err,%d!=%d,body:%s", n, p.bodylength, string(body))
			return body[:n], nil
		}
	}
	return body, nil
}

// Connection Wrapper around net.Conn.
type Connection interface {
	net.Conn
	Network() string
	// String() string
	ReadFrom(buf []byte) (num int, raddr net.Addr, err error)
	WriteTo(buf []byte, raddr net.Addr) (num int, err error)
}

// Connection implementation.
type connection struct {
	baseConn net.Conn
	laddr    net.Addr
	raddr    net.Addr
	// mu       sync.RWMutex
	logKey string
}

func newUDPConnection(baseConn net.Conn) Connection {
	conn := &connection{
		baseConn: baseConn,
		laddr:    baseConn.LocalAddr(),
		raddr:    baseConn.RemoteAddr(),
		logKey:   "udpConnection",
	}
	return conn
}

func (conn *connection) Read(buf []byte) (int, error) {
	var (
		num int
		err error
	)

	num, err = conn.baseConn.Read(buf)

	if err != nil {
		return num, utils.NewError(err, conn.logKey, "read", conn.baseConn.LocalAddr().String())
	}
	return num, err
}

func (conn *connection) ReadFrom(buf []byte) (num int, raddr net.Addr, err error) {
	num, raddr, err = conn.baseConn.(net.PacketConn).ReadFrom(buf)
	if err != nil {
		return num, raddr, utils.NewError(err, conn.logKey, "readfrom", conn.baseConn.LocalAddr().String(), raddr.String())
	}
	//logrus.Tracef("readFrom %d , %s -> %s \n %s", num, raddr, conn.LocalAddr(), string(buf[:num]))
	return num, raddr, err
}

func (conn *connection) Write(buf []byte) (int, error) {
	var (
		num int
		err error
	)

	num, err = conn.baseConn.Write(buf)
	if err != nil {
		return num, utils.NewError(err, conn.logKey, "write", conn.baseConn.LocalAddr().String())
	}
	return num, err
}

func (conn *connection) WriteTo(buf []byte, raddr net.Addr) (num int, err error) {
	num, err = conn.baseConn.(net.PacketConn).WriteTo(buf, raddr)
	if err != nil {
		return num, utils.NewError(err, conn.logKey, "writeTo", conn.baseConn.LocalAddr().String(), raddr.String())
	}
	//logrus.Tracef("writeTo %d , %s -> %s \n %s", num, conn.baseConn.LocalAddr(), raddr.String(), string(buf[:num]))
	return num, err
}

func (conn *connection) LocalAddr() net.Addr {
	return conn.baseConn.LocalAddr()
}

func (conn *connection) RemoteAddr() net.Addr {
	return conn.baseConn.RemoteAddr()
}

func (conn *connection) Close() error {
	err := conn.baseConn.Close()
	if err != nil {
		return utils.NewError(err, conn.logKey, "close", conn.baseConn.LocalAddr().String(), conn.baseConn.RemoteAddr().String())
	}
	return nil
}

func (conn *connection) Network() string {
	return strings.ToUpper(conn.baseConn.LocalAddr().Network())
}

func (conn *connection) SetDeadline(t time.Time) error {
	return conn.baseConn.SetDeadline(t)
}

func (conn *connection) SetReadDeadline(t time.Time) error {
	return conn.baseConn.SetReadDeadline(t)
}

func (conn *connection) SetWriteDeadline(t time.Time) error {
	return conn.baseConn.SetWriteDeadline(t)
}
