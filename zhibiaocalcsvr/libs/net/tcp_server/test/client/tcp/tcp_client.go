package tcp

import (
	"fmt"
	"math/rand"
	"net"
	"libs/net/tcp_server/codec"
	"time"
)

type TcpClient struct {
	codec        *codec.Codec
	Name         string
}

type BHDR struct {
	//UnLen  uint16
	Cmd  uint32
	Seq  uint32
	Type uint8
	Uid  uint64
}

var r *rand.Rand

func init(){
	r = rand.New(rand.NewSource(time.Now().Unix()))
}

func (c *TcpClient) Start() {
	conn, err := net.Dial("tcp", "127.0.0.1:50000")
	if err != nil {
		fmt.Println(err)
		return
	}

	c.codec = codec.NewCodec(conn)

	go c.HeartBeat()
	go c.Receive()
}



func (c *TcpClient) HeartBeat() {
	ticker := time.NewTicker(time.Minute * 4)
	for _ = range ticker.C {
		fmt.Println("heart beat...")
		msgHead := common.BHDR{
			Cmd:  1,
			Seq:  2,
			Type: 44,
			Uid: r.Uint64(),
		}
		err := c.codec.Eecode(codec.Package{PackageHead: msgHead, Content: []byte{}}, 10*time.Second)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (c *TcpClient) Receive() {
	for {
		_, err := c.codec.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		for {
			fmt.Println()
			pack, ok := c.codec.Decode()
			if ok {
				c.HandlePackage(*pack)
				continue
			}
			break
		}
	}
}

func (c *TcpClient) HandlePackage(pack codec.Package) error {

	return c.handleMessage(pack)
}


func (c *TcpClient) handleMessage(pack codec.Package) error {
	fmt.Println(string(pack.Content))
	return nil
}

func (c *TcpClient) SendMessage() {
	msgHead := common.BHDR{
		Cmd:  1,
		Seq:  2,
		Type: 3,
		Uid: r.Uint64(),
	}
	err := c.codec.Eecode(codec.Package{PackageHead: msgHead, Content: []byte(c.Name)}, 10*time.Second)
	if err != nil {
		fmt.Println(err)
	}
}
