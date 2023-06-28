package customError

import (
	"fmt"
	"runtime/debug"
)

type MyError struct {
	Inner      error // 用于存储要封装的异常
	Message    string
	StackTrace string                 // 用于记录当异常发生时的堆栈跟踪信息
	Misc       map[string]interface{} // 用于存储其他杂项信息
}

func WrapError(err error, formatMsg string, msgArgs ...interface{}) MyError {
	return MyError{
		Inner:      err,
		Message:    fmt.Sprintf(formatMsg, msgArgs...),
		StackTrace: string(debug.Stack()),
		Misc:       make(map[string]interface{}),
	}
}

func (e MyError) Error() string {
	return e.Message
}
