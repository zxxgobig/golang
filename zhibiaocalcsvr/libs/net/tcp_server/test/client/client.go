package main

import (
	"libs/net/tcp_server/test/client/tcp"
	"time"
)

func main() {
	TestClient()
}

func TestClient() {
	client1 := tcp.TcpClient{Name:"c1"}
	client2 := tcp.TcpClient{Name:"c2"}
	client1.Start()
	client2.Start()


	go func() {
		for {
			time.Sleep(1*time.Second)
			client1.SendMessage()
		}
	}()


	for {
		time.Sleep(1*time.Second)
		client2.SendMessage()
	}


}
