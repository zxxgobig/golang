package main

import (
	"libs/log"
	"libs/net/tcp_server/codec"
	"libs/net/tcp_server/connect"
)

func main() {
	// 启动服务器
	conf := connect.Conf{
		Address:      "127.0.0.1" + ":" + "50000",
		MaxConnCount: 100,
		AcceptCount:  1,
	}

	server := connect.NewTCPServer(conf)
	server.RegisterOnTcpConnct(OnTcpConnect)
	server.RegisterOnTcpReceive(OnTcpReceive)
	server.RegisterOnTcpClose(OnTcpClose)

	server.Start()
}

func OnTcpConnect(conn_ctx *connect.ConnContext) {
	log.Release("tcp connect [%s] on create,id[%d]", conn_ctx.GetRemoteAddr(), conn_ctx.GetId())
}

func OnTcpReceive(conn_ctx *connect.ConnContext, recivePack *codec.Package) {
	log.Release("receive msg from tcp connect [%s]", conn_ctx.GetRemoteAddr)
	log.Release("receive package: %+v", *recivePack)

	conn_ctx.SendMessage(&codec.Package{PackageHead: recivePack.PackageHead, Content: []byte("msg from server")})
}

func OnTcpClose(conn_ctx *connect.ConnContext) {
	log.Release("tcp connect [%s] on closed", conn_ctx.GetRemoteAddr)
}
