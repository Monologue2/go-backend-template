package routes

import (
	"github.com/gin-gonic/gin"
	test "shogle.net/template/api/handlers/test"
)

func SetupTestRoutes(r *gin.Engine, testHandler *test.TestHandler) {
	testGroup := r.Group("/api/test")
	{
		testGroup.GET("/ping", test.Ping)
		testGroup.POST("/", testHandler.CreateTest)
		testGroup.GET("/", testHandler.GetTests)
		testGroup.GET("/:id", testHandler.GetTestByID)
		testGroup.DELETE("/:id", testHandler.DeleteTest)
	}
}
