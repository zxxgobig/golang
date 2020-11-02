package codec

import (
	"bytes"
	"encoding/binary"
	"github.com/pkg/errors"
	"libs/log"
	"libs/net/tcp_server/common/config"
	"net"
	"reflect"
	"time"
)

type BHDR struct {
	//UnLen  uint16
	Cmd  uint32
	Seq  uint32
	Type uint8
	Uid  uint64
}

type Codec struct {
	Conn       net.Conn
	MsgHeadLen int
	ReadBuf    buffer // 读缓冲
	WriteBuf   []byte // 写缓冲
}

// newCodec 创建一个解码器
func NewCodec(conn net.Conn) *Codec {
	return &Codec{
		Conn:       conn,
		ReadBuf:    newBuffer(conn, config.MaxRecvMshLen+4),
		WriteBuf:   make([]byte, config.MaxSendMshLen+4),
		MsgHeadLen: binary.Size(BHDR{}),
	}
}

// Read 从conn里面读取数据，当conn发生阻塞，这个方法也会阻塞
func (c *Codec) Read() (int, error) {
	return c.ReadBuf.readFromReader()
}

// Decode 解码数据
func (c *Codec) Decode() (*Package, bool) {
	var err error

	// 读取包长
	bufMsgLen, err := c.ReadBuf.seek(0, config.LenMsgLen)
	if err != nil {
		return nil, false
	}

	// parse len
	var msgLen uint32
	switch config.LenMsgLen {
	case 1:
		msgLen = uint32(bufMsgLen[0])
	case 2:
		if config.Binary_endian == binary.LittleEndian {
			msgLen = uint32(binary.LittleEndian.Uint16(bufMsgLen))
		} else {
			msgLen = uint32(binary.BigEndian.Uint16(bufMsgLen))
		}
	case 4:
		if config.Binary_endian == binary.LittleEndian {
			msgLen = binary.LittleEndian.Uint32(bufMsgLen)
		} else {
			msgLen = binary.BigEndian.Uint32(bufMsgLen)
		}
	}

	if config.MsgLenIncludeHead {
		msgLen -= uint32(config.LenMsgLen)
	}

	if msgLen > config.MaxRecvMshLen {
		return nil, false
	}

	if msgLen > 0 {
		headOffset := config.LenMsgLen + c.MsgHeadLen
		msghead, err := c.ReadBuf.seek(config.LenMsgLen, headOffset)
		if err != nil {
			log.Error("read package head fail: %s", err)
			return nil, false
		}
		messageHead := &BHDR{}
		buffer := new(bytes.Buffer)
		buffer.Write(msghead)
		err = binary.Read(buffer, config.Binary_endian, messageHead)
		if err != nil {
			log.Error("package head read by Binary_endian fail: %s", err)
			return nil, false
		}

		limit := int(msgLen) - c.MsgHeadLen
		msgBody, err := c.ReadBuf.read(headOffset, limit)
		if err != nil {
			log.Error("read package body fail: %s", err)
			return nil, false
		}

		message := Package{PackageHead: *messageHead, Content: msgBody}
		return &message, true

	} else {
		return nil, false
	}

}

// Eecode 编码数据
func (c *Codec) Eecode(pack Package, duration time.Duration) error {
	if reflect.DeepEqual(pack.PackageHead, BHDR{}) {
		return errors.New("send data has not BHDR head")
	}

	var msgLen uint32
	buf := bytes.NewBuffer(nil)
	msgHead := pack.PackageHead
	msgBody := pack.Content
	binary.Write(buf, config.Binary_endian, msgHead)
	binary.Write(buf, config.Binary_endian, msgBody)
	msgBytes := buf.Bytes()
	msgLen += uint32(len(msgBytes))

	if msgLen > config.MaxSendMshLen {
		return errors.New("write message too long")
	}

	msg := make([]byte, uint32(config.LenMsgLen)+msgLen)

	if config.MsgLenIncludeHead {
		msgLen += uint32(config.LenMsgLen)
	}

	// write len
	switch config.LenMsgLen {
	case 1:
		msg[0] = byte(msgLen)
	case 2:
		if config.Binary_endian == binary.LittleEndian {
			binary.LittleEndian.PutUint16(msg, uint16(msgLen))
		} else {
			binary.BigEndian.PutUint16(msg, uint16(msgLen))
		}
	case 4:
		if config.Binary_endian == binary.LittleEndian {
			binary.LittleEndian.PutUint32(msg, msgLen)
		} else {
			binary.BigEndian.PutUint32(msg, msgLen)
		}
	}

	// write data
	l := config.LenMsgLen

	copy(msg[l:], msgBytes)

	c.Conn.SetWriteDeadline(time.Now().Add(duration))
	_, err := c.Conn.Write(msg)
	if err != nil {
		return err
	}
	return nil
}
