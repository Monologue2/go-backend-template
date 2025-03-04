package routes

import (
	test "project/api/handlers/test"

	"github.com/gin-gonic/gin"
)

func SetupTestRoutes(r *gin.Engine, testHandler *test.TestHandler) {
	testGroup := r.Group("/api/test")
	{
		testGroup.GET("/ping", test.Ping)
		testGroup.POST("/", testHandler.CreateTest)
		testGroup.GET("/", testHandler.GetTests)
		testGroup.GET("/:id", testHandler.GetTestByID)
		testGroup.DELETE("/:id", testHandler.DeleteTest)
		testGroup.GET("/add/:id", testHandler.AddOne)
	}
}
