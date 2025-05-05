package internal

import (
	"TZ/internal/handler"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.GET("/", handler.GetallPerson)
	api.POST("/people", handler.CreatePerson)
	api.GET("/people/:id", handler.GetPersonByID)
	api.DELETE("/people/:id", handler.DeletePerson)
	api.PUT("/people/:id", handler.UpdatePerson)
	api.PATCH("/people/:id", handler.PatchPerson)
	r.StaticFile("/swagger/openapi.yaml", "./docs/openapi.yaml")

}
