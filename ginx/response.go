package ginx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func ResponseOk(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{Code: 0, Message: "ok", Data: data})
}

func ResponseFail(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{Code: code, Message: message})
}

func ResponseError(c *gin.Context, err error) {
	if err == nil {
		ResponseOk(c, nil)
		return
	}
	if e, ok := err.(*Error); ok {
		c.JSON(http.StatusOK, Response{Code: e.Code, Message: e.Message})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 500, Message: err.Error()})
}
