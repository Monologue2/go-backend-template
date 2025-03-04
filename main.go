package main

import (
	test "project/api/handlers/test"

	"project/api/routes"
	"project/api/services"
	"project/middlewares"
	"project/models"
	"project/repositories"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), middlewares.RequestLogger())

	models.Migrate(repositories.DB)
	userService := services.NewTestService(repositories.DB)
	userHandler := test.NewTestHandler(userService)

	// :8080/api/test/ping {"message":"pong"}
	routes.SetupTestRoutes(r, userHandler)
	r.Run(":8080")
}
