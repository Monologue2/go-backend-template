package routes

import (
	"github.com/gin-gonic/gin"
	test "shogle.net/template/handlers/test"
)

func SetupTestRoutes(r *gin.Engine) {
	ping := test.Ping

	testGroup := r.Group("/api/test")
	{
		testGroup.GET("/ping", ping)
	}
}
