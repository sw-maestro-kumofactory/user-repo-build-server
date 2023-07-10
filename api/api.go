package api

import (
	"net/http"
	"github.com/gin-gonic/gin"

	"github.com/sw-maestro-kumofactory/miz-ball/api/v1.0"
)

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

// ApplyRoutes applies router to gin Router
func ApplyRoutes(r *gin.Engine) {
	r.GET("/", health)
	api := r.Group("/api")
	{
		apiv1.ApplyRoutes(api)
	}
}