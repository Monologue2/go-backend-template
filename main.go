package main

import (
	"github.com/gin-gonic/gin"
	test "shogle.net/template/api/handlers/test"
	"shogle.net/template/api/routes"
	"shogle.net/template/api/services"
	"shogle.net/template/middlewares"
	"shogle.net/template/models"
	"shogle.net/template/repositories"
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
