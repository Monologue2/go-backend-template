package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		log.Printf("[%s] %s | %s | %dms",
			c.Request.Method, c.Request.URL.Path, c.ClientIP(), duration.Microseconds())
	}
}
