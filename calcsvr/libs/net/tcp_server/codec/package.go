package codec

// Package 消息包
type Package struct {
	PackageHead BHDR   // 固定包头
	Content     []byte // 消息体(proto)
}
