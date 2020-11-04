package global

type IZhiBiao interface {
	Start() error
	Stop()
	GetIndexCodes() map[string]struct{}
}
