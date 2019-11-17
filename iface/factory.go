package iface

// 日志工厂接口
type Factory interface {
	GetDefaultLogger() Logger
	GetLogger(name string) Logger
}
