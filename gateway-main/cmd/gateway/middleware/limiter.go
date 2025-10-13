package middleware

import (
	"github.com/BitofferHub/gateway/internal/conf"
	limit "github.com/BitofferHub/gateway/limiter"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"net/http"
)

// 限流器
func Limiter() gin.HandlerFunc {
	return limiter
}

func limiter(c *gin.Context) {
	action, _ := c.Params.Get("action")
	if _, ok := conf.Routes[action]; !ok {
		c.Abort()
		c.JSON(http.StatusNotFound, "route err")
		return
	}
	result, err := limit.Rl.Allow(c.Request.Context(), action)

	if err != nil {
		log.Errorf("rate limiter error %+v\n", err)
		c.Abort()
		c.JSON(http.StatusOK, "ok")
		return
	}

	if result.IsAllowed {
		c.Next()
		return
	} else {
		c.Abort()
		c.JSON(http.StatusOK, "ok")
		return
	}

}
