package apiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	// "github.com/sw-maestro-kumofactory/miz-ball/api/v1.0/{api_name}"
	"github.com/sw-maestro-kumofactory/miz-ball/api/v1.0/sample"
)

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// ApplyRoutes applies router to gin Router
func ApplyRoutes(r *gin.RouterGroup) {
	v1 := r.Group("/v1")
	{
		v1.GET("/ping", ping)
		v1.GET("/TEST_CREATE", sample.SAMPLE_TEST_CREATE)
		v1.GET("/TEST_BUILD", sample.SAMPLE_TEST_CREATE)
	}
}
