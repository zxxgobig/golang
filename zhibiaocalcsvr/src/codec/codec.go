package codec

import (
	"bytes"
	"common"
	"encoding/binary"
	"errors"
	"github.com/golang/protobuf/proto"
	"libs/net/tcp"
)

type ClientCodec struct {
}

func (this *ClientCodec) NewCodec() tcp.Codec {
	return &BinaryCodec{
		msg_head_len: binary.Size(common.BHDR{}),
	}
}

type BinaryCodec struct {
	msg_head_len int
}

//编码包头和包体成消息
func (bc *BinaryCodec) Marshal(v interface{}) ([]byte, error) {
	msg, ok := v.([]interface{})
	if !ok {
		return nil, errors.New("invalid send data")
	}

	msg_head, ok := msg[0].(*common.BHDR)
	if !ok {
		return nil, errors.New("send data has not BHDR head")
	}

	buf := bytes.NewBuffer(nil)
	binary.Write(buf, binary.BigEndian, msg_head)
	if len(msg) == 2 {
		var body_data []byte
		var err error
		body, ok := msg[1].(proto.Message)
		if ok {
			body_data, err = proto.Marshal(body)
			if err != nil {
				return nil, err
			}
		} else {
			body_data, ok = msg[1].([]byte)
			if !ok {
				return nil, errors.New("send data protomsg must be proto or []byte")
			}
		}
		binary.Write(buf, binary.BigEndian, body_data)
	}

	return buf.Bytes(), nil
}

//解析除消息长度之外的字节，解析为消息头和消息体
func (bc *BinaryCodec) Unmarshal(data []byte) (interface{}, error) {
	if len(data) < bc.msg_head_len {
		return nil, errors.New("invalid length")
	}
	msg := []interface{}{}
	//解码成一个消息结构体
	msg_head := &common.BHDR{}

	buf_data := bytes.NewBuffer(data)
	binary.Read(buf_data, binary.BigEndian, msg_head)

	msg = append(msg, msg_head)

	//包头之后是具体的消息内容
	msg_body := data[bc.msg_head_len:]
	msg = append(msg, msg_body)

	return msg, nil
}
