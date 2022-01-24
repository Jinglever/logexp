package logexp

import (
	"encoding/json"
	"fmt"
)

// Custom Error
type CstError struct {
	Code    int    `json:"code"` //code 错误码
	Message string `json:"msg"`  //msg 消息
}

func (err *CstError) Error() string {
	b, _ := json.Marshal(err)
	return string(b)
}

var (
	ErrCodeUnknown           = 10001 // not sure exactly the error meaning
	ErrCodeInvalidExpression = 10002 // invalid expression
)

func newCstError(code int, format string, a ...interface{}) *CstError {
	return &CstError{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
	}
}
