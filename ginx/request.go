package ginx

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func BindForm(c *gin.Context, obj any) bool {
	if err := c.ShouldBind(obj); err != nil {
		ResponseFail(c, http.StatusBadRequest, err.Error())
		return false
	}
	return true
}

func BindUri(c *gin.Context, obj any) bool {
	if err := c.ShouldBindUri(obj); err != nil {
		ResponseFail(c, http.StatusBadRequest, err.Error())
		return false
	}
	return true
}

func BindJSON(c *gin.Context, obj any) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		ResponseFail(c, http.StatusBadRequest, err.Error())
		return false
	}
	return true
}

func BindQuery(c *gin.Context, obj any) bool {
	if err := c.ShouldBindQuery(obj); err != nil {
		ResponseFail(c, http.StatusBadRequest, err.Error())
		return false
	}
	return true
}

func QueryString(c *gin.Context, key string) string {
	return c.Query(key)
}

func QueryInt(c *gin.Context, key string) int {
	return cast.ToInt(c.Query(key))
}

func QueryFloat(c *gin.Context, key string) float64 {
	return cast.ToFloat64(c.Query(key))
}

func QueryBool(c *gin.Context, key string) bool {
	return cast.ToBool(c.Query(key))
}
