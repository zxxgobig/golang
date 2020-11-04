package util

// 通过类型断言判断具体的网络错误，替代项目中使用字符串匹配方式的错误类型判断
// 不可用，待完成

//import (
//	"log"
//	"net"
//	"os"
//	"syscall"
//)
//
//
//func IsClientCloseConn(err error) bool {
//	netErr, ok := err.(net.Error)
//	if !ok {
//		return false
//	}
//
//	opErr, ok := netErr.(*net.OpError)
//	if !ok {
//		return false
//	}
//
//	switch t := opErr.Err.(type) {
//	case *os.SyscallError:
//		if errno, ok := t.Err.(syscall.Errno); ok {
//			switch errno {
//			case syscall.ECONNRESET:
//				log.Println("连接重置")
//				return true
//			}
//		}
//	}
//	return false
//}
//
//func IsServerCloseConn(err error) bool {
//	netErr, ok := err.(net.Error)
//	if !ok {
//		return false
//	}
//
//	opErr, ok := netErr.(*net.OpError)
//	if !ok {
//		return false
//	}
//
//	switch t := opErr.Err.(type) {
//	case *os.SyscallError:
//		if errno, ok := t.Err.(syscall.Errno); ok {
//			switch errno {
//			case syscall.ECONNRESET:
//				log.Println("连接重置")
//				return true
//			}
//		}
//	}
//	return false
//}
//
//func IsReadTimeOut(err error) bool {
//	netErr, ok := err.(net.Error)
//	if !ok {
//		return false
//	}
//	return netErr.Timeout()
//}
