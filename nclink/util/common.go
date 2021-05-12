package util

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"
)

func GoSafely(handler func(), catchFunc func(r interface{})) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr, "%s goroutine panic: %v\n%s\n",
					time.Now(), r, string(debug.Stack()))
				if catchFunc != nil {
					go func() {
						defer func() {
							if p := recover(); p != nil {
								fmt.Fprintf(os.Stderr, "recover goroutine panic:%v\n%s\n",
									p, string(debug.Stack()))
							}
						}()
						catchFunc(r)
					}()
				}
			}
		}()
		handler()
	}()
}

func TimeToUnixMs(t time.Time) int64 {
	return t.Unix()*1000 + int64(t.Nanosecond()/10000000)
}

/*
func ErrorWrap(err error, msg string) error {
	logger.Errorf("error:%v msg:%s\n%+v", err, msg, callers())
	return err
}

func NewError(msg string, format string, args ...interface{}) error {
	err := fmt.Errorf(format, args...)
	logger.Errorf("error:%v msg:%s\n%+v", err, msg, callers())
	return err
}

func callers() []uintptr {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])
	d := pcs[0:n]
	return d
}
*/
