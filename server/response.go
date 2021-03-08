//Created by Goland
//@User: lenora
//@Date: 2021/1/15
//@Time: 10:34 上午

package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ApiResponse struct {
	Code    int16       `json:"code"`    // 状态码
	Message string      `json:"message"` // 状态短语
	Result  interface{} `json:"result"`  // 数据结果集
}

const (
	SUCCESS int16 = 1000
	FAIL    int16 = 4000

	NOT_FOUND int16 = 3001

	WRONG_PARAM int16 = 4001
)

var WrongMessage = map[int16]string{
	1000: "success",

	2001: "user licence expired",
	2002: "user wrong",
	2003: "permission denied",

	3001: "record not found",
	3002: "record has been exist",

	4000: "fail",
	4001: "param error",

	5001: "connect source error",
}

func (this *defaultServer) InvalidParametersError(c *gin.Context) {
	responseError(c, WRONG_PARAM, nil)
}

func (this *defaultServer) InternalServiceError(c *gin.Context, message ...string) {
	responseError(c, FAIL, message)
}

func (this *defaultServer) ResponseError(c *gin.Context, code int16, message ...string) {
	responseError(c, code, message)
}

func (this *defaultServer) ResponseSuccess(c *gin.Context, result interface{}, msg ...string) {
	responseOutput(c, SUCCESS, msg, result)
}

func responseError(c *gin.Context, code int16, message []string) {
	responseOutput(c, code, message, nil)
}

func responseOutput(c *gin.Context, code int16, message []string, result interface{}) {
	var msg string
	if len(message) <= 0 {
		msg = WrongMessage[code]
	} else {
		msg = message[0]
	}
	c.JSON(http.StatusOK, ApiResponse{
		Code:    code,
		Message: msg,
		Result:  result,
	})
}
