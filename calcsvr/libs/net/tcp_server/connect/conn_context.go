package connect

import (
	"fmt"
	"io"
	"libs/net/tcp_server/codec"
	"libs/log"
	"net"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

const (
	ReadDeadline  = 10 * time.Minute
	WriteDeadline = 10 * time.Second
)

var global_id uint64

type DealFunc func(*codec.Package, chan *codec.Package)

// ConnContext 连接上下文
type ConnContext struct {
	id      uint64
	stats   bool
	Codec   *codec.Codec // 编解码器
	msgSend chan *codec.Package
}

func NewConnContext(conn *net.TCPConn) *ConnContext {
	cc := codec.NewCodec(conn)
	mc := make(chan *codec.Package, 100)
	current_id := atomic.AddUint64(&global_id, 1)
	return &ConnContext{Codec: cc, msgSend: mc, stats: true, id: current_id}
}

// DoConn 处理TCP连接
func (c *ConnContext) DoConn(t *TCPServer) {
	defer RecoverPanic()

	go func() {
	Loop:
		for {
			select {
			case msg := <-c.msgSend:
				if msg != nil {
					c.sendMessage(msg)
				}

				if !c.stats {
					break Loop
				}
			}
		}
	}()

	for {
		err := c.Codec.Conn.SetReadDeadline(time.Now().Add(ReadDeadline))
		if err != nil {
			c.HandleReadErr(err)
			return
		}

		_, err = c.Codec.Read()
		if err != nil {
			c.HandleReadErr(err)
			break
		}

		for {
			message, ok := c.Codec.Decode()
			if ok {
				func() {
					defer func() {
						if err := recover(); err != nil {
							log.Release("on tcp receive panic: %s", err)
							buf := make([]byte, 2048)
							n := runtime.Stack(buf, false)
							errinfo := fmt.Sprintf("%s", buf[:n])
							log.Error("on tcp receive panic: %s", errinfo)
						}
					}()
					t.onTcpReceive(c, message)
				}()
				continue
			}
			break
		}
	}

	func() {
		defer func() {
			if err := recover(); err != nil {
				log.Release("on tcp close panic: %s", err)
				buf := make([]byte, 2048)
				n := runtime.Stack(buf, false)
				errinfo := fmt.Sprintf("%s", buf[:n])
				log.Error("on tcp close panic: %s", errinfo)
			}
		}()
		t.onTcpClose(c)
	}()

}

func (c *ConnContext) GetId() uint64 {
	return c.id
}

func (c *ConnContext) GetRemoteAddr() net.Addr {
	return c.Codec.Conn.RemoteAddr()
}

func (c *ConnContext) SendMessage(p *codec.Package) {
	if c.stats && p != nil {
		c.msgSend <- p
	}
}

func (c *ConnContext) sendMessage(p *codec.Package) {
	err := c.Codec.Eecode(*p, WriteDeadline)
	if err != nil {
		log.Error("send message to client fail: %s", err)
	} else {
		log.Release("send message to client success")
	}
}

// HandleReadErr 读取conn错误
func (c *ConnContext) HandleReadErr(err error) {
	str := err.Error()
	log.Release("connect err: %s", err)

	// 服务器主动关闭连接
	if strings.HasSuffix(str, "use of closed network connection") {
		return
	}

	c.Release()

	// 客户端主动关闭连接或者异常程序退出
	if err == io.EOF {
		return
	}
	// SetReadDeadline 之后，超时返回的错误
	if strings.HasSuffix(str, "i/o timeout") {
		return
	}
}

// Release 释放TCP连接
func (c *ConnContext) Release() {
	err := c.Codec.Conn.Close()
	if err != nil {
		log.Error("connect close fail: %s", err)
	}
	close(c.msgSend)
	c.stats = false
}
