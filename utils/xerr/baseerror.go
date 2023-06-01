package xerr

import (
	"fmt"
)

// BusinessError 常用通用固定错误
type BusinessError struct {
	errCode uint32
	errMsg  string
}

type Option func(*BusinessError)

// GetErrCode 返回给前端的错误码
func (e *BusinessError) GetErrCode() uint32 {
	return e.errCode
}

// GetErrMsg 返回给前端显示端错误信息
func (e *BusinessError) GetErrMsg() string {
	return e.errMsg
}

func (e *BusinessError) Error() string {
	return fmt.Sprintf("BusinessError:%d，ErrMsg:%s", e.errCode, e.errMsg)
}

func SetCode(code uint32) Option { // 返回一个Option类型的函数（闭包）：接受BusinessError类型指针参数并修改
	return func(e *BusinessError) {
		e.errCode = code
	}
}

func SetMsg(msg string) Option { // 返回一个Option类型的函数（闭包）：接受BusinessError类型指针参数并修改
	return func(e *BusinessError) {
		e.errMsg = msg
	}
}

func NewBusinessError(ops ...Option) *BusinessError {
	//初始化默认值
	defaultError := BusinessError{
		errCode: ErrorBusiness,
		errMsg:  "",
	}
	//依次调用opts函数列表中的函数，为结构体成员赋值
	for _, opt := range ops {
		opt(&defaultError)
	}
	return &defaultError
}