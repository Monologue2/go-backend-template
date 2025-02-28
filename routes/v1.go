package routes

import (
	"github.com/gin-gonic/gin"
	"shogle.net/template/handlers/v1"
)

func SetupV1Routes(r *gin.Engine) {
	ping := handlers.Ping

	v1Group := r.Group("/api/v1")
	{
		v1Group.GET("/ping", ping)
	}
}
