package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

const (
	BAD_REQUEST     = 400
	ERROR           = 500
	SUCCESS         = 200
	UNAUTHORIZATION = 401
)

func SuccessResult(code int, data interface{}, msg string, c *gin.Context) {
	c.IndentedJSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

func ValidationFailureResult(code int, data interface{}, msg string, c *gin.Context) {
	c.IndentedJSON(http.StatusBadRequest, Response{
		code,
		data,
		msg,
	})
}

func InternalErrorResult(code int, data interface{}, msg string, c *gin.Context) {
	c.IndentedJSON(http.StatusInternalServerError, Response{
		code,
		data,
		msg,
	})
}

func NoAuth(message string, c *gin.Context) {
	c.IndentedJSON(http.StatusUnauthorized, Response{
		UNAUTHORIZATION,
		nil,
		message,
	})
}

func Ok(c *gin.Context) {
	SuccessResult(SUCCESS, map[string]interface{}{}, "success", c)
}

func OkWithMessage(message string, c *gin.Context) {
	SuccessResult(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	SuccessResult(SUCCESS, data, "success", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	SuccessResult(SUCCESS, data, message, c)
}

func ValidationFailureWithMessage(message string, c *gin.Context) {
	ValidationFailureResult(BAD_REQUEST, map[string]interface{}{}, message, c)
}

func InteralErrorWithMessage(message string, c *gin.Context) {
	ValidationFailureResult(ERROR, map[string]interface{}{}, message, c)
}
