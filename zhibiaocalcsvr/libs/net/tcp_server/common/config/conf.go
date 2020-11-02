package config

import (
	"encoding/binary"
)


var (
	Binary_endian binary.ByteOrder = binary.BigEndian //表示大端还是小端
	MaxRecvMshLen uint32 = 2048 * 1024
	MaxSendMshLen uint32 = 2048 * 1024
	LenMsgLen =  4   //协议包长字段长度
	MsgLenIncludeHead = true  //包长是否包含
)



