package connect

import (
	"fmt"
	"libs/log"
	"libs/net/tcp_server/codec"
	"net"
	"runtime"
)

// Conf server配置文件
type Conf struct {
	Address string //
	// 端口
	MaxConnCount int // 最大连接数
	AcceptCount  int // 接收建立连接的groutine数量
}

// TCPServer TCP服务器
type TCPServer struct {
	Address      string // 端口
	MaxConnCount int    // 最大连接数
	AcceptCount  int    // 接收建立连接的groutine数量
	onTcpConnct  func(ctx *ConnContext)
	onTcpReceive func(ctx *ConnContext, p *codec.Package)
	onTcpClose   func(ctx *ConnContext)
}

// NewTCPServer 创建TCP服务器
func NewTCPServer(conf Conf) *TCPServer {
	return &TCPServer{
		Address:      conf.Address,
		MaxConnCount: conf.MaxConnCount,
		AcceptCount:  conf.AcceptCount,
	}
}

func (t *TCPServer) RegisterOnTcpConnct(f func(ctx *ConnContext)) {
	t.onTcpConnct = f
}

func (t *TCPServer) RegisterOnTcpReceive(f func(ctx *ConnContext, rec_pack *codec.Package)) {
	t.onTcpReceive = f
}

func (t *TCPServer) RegisterOnTcpClose(f func(ctx *ConnContext)) {
	t.onTcpClose = f
}

// Start 启动服务器
func (t *TCPServer) Start() {
	addr, err := net.ResolveTCPAddr("tcp", t.Address)
	if err != nil {
		log.Error("Server start error : %s", err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Error("error listening: %s", err.Error())
		return
	} else {
		log.Release("server listen at %s", addr)
	}
	for i := 0; i < t.AcceptCount; i++ {
		go t.Accept(listener)
	}

	select {}
}

// Accept 接收客户端的TCP长连接的建立
func (t *TCPServer) Accept(listener *net.TCPListener) {
	defer RecoverPanic()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Error("accept tcp err: %s", err)
			continue
		}

		err = conn.SetKeepAlive(true)
		if err != nil {
			log.Error("set conn keep live err: %s", err)
		}

		connContext := NewConnContext(conn)

		func() {
			defer func() {
				if err := recover(); err != nil {
					log.Release("on tcp connect panic: %s", err)
					buf := make([]byte, 2048)
					n := runtime.Stack(buf, false)
					errinfo := fmt.Sprintf("%s", buf[:n])
					log.Release("on tcp connect panic: %s", errinfo)
				}
			}()
			t.onTcpConnct(connContext)
		}()

		go connContext.DoConn(t)
	}
}

// RecoverPanic 恢复panic
func RecoverPanic() {
	err := recover()
	if err != nil {
		log.Error(GetPanicInfo())
	}

}

// PrintStaStack 打印Panic堆栈信息
func GetPanicInfo() string {
	buf := make([]byte, 2048)
	n := runtime.Stack(buf, false)
	return fmt.Sprintf("%s", buf[:n])
}
