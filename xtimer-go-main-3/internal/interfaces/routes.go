// routes.go

package interfaces

import (
	"bytes"
	"context"
	"io/ioutil"
	"time"

	"github.com/BitofferHub/pkg/constant"
	engine "github.com/BitofferHub/pkg/middlewares/gin"
	"github.com/BitofferHub/pkg/middlewares/log"
	"github.com/BitofferHub/xtimer/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	xTimerService *service.XTimerService
}

func NewHandler(s *service.XTimerService) *Handler {
	return &Handler{
		xTimerService: s,
	}
}

func NewRouter(h *Handler) *gin.Engine {
	r := engine.NewEngine(engine.WithLogger(false))
	r.Use(Cors())
	// 使用gin中间件
	r.Use(InfoLog())
	project := r.Group("xtimer")
	project.POST("/createTimer", h.CreateTimer)
	project.GET("/enableTimer", h.EnableTimer)
	project.GET("/unableTimer", h.UnableTimer)
	project.DELETE("/delTimer", h.DelTimer)
	project.GET("/getTimers", h.GetTimers)
	project.POST("/callback", h.TestCallback)
	return r
}

func InfoLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		beginTime := time.Now()
		// ***** 1. get request body ****** //
		traceID := c.Request.Header.Get(constant.TraceID)
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body.Close() //  must close
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		// ***** 2. set requestID for goroutine ctx ****** //
		// duration := float64(time.Since(beginTime)) / float64(time.Second)
		ctx := context.WithValue(context.Background(), constant.TraceID, traceID)
		log.InfoContextf(ctx, "ReqPath[%s]-Cost[%v]\n", c.Request.URL.Path, time.Since(beginTime))
	}
}
