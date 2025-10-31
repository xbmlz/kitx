package ginx

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/requestid"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/xbmlz/kitx/log"
)

type Engine struct {
	*gin.Engine
}

func New() *Engine {
	gin.SetMode(gin.ReleaseMode)

	logger := log.GetLogger()

	r := gin.New()
	r.Use(ginzap.Ginzap(logger, time.DateTime, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))
	r.Use(cors.Default())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(requestid.New())

	r.GET("/healthz", func(ctx *gin.Context) {
		reqID := requestid.Get(ctx)
		ctx.JSON(http.StatusOK, gin.H{"status": "ok", "request_id": reqID, "time": time.Now().Format(time.DateTime)})
	})

	return &Engine{r}
}
