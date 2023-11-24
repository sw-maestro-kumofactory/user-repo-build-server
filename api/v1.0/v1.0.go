package apiv1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/sw-maestro-kumofactory/miz-ball/api/v1.0/deploy"
	"github.com/sw-maestro-kumofactory/miz-ball/api/v1.0/rds"
	"github.com/sw-maestro-kumofactory/miz-ball/middleware"
)

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func ApplyRoutes(r *gin.RouterGroup) {
	v1 := r.Group("/v1")
	{
		v1.GET("/ping", ping)
		// v1.GET("/TEST_CREATE", sample.SAMPLE_TEST_CREATE)
		// v1.GET("/TEST_BUILD", sample.SAMPLE_TEST_BUILD2)
		// v1.GET("/TEST_CLONE", sample.SAMPLE_TEST_CLONE)
		// v1.GET("/TEST_PUSH", sample.SAMPLE_TEST_PUSH)
		// v1.GET("/TEST_FIND", sample.FindDockerfile)
		// v1.GET("/TEST_EXTRACT", sample.SAMPLE_EXTRACT)
		// v1.GET("/TEST_ARCHIVE", sample.SAMPLE_ARCHIVE)
		v1.POST("/deploy", deploy.ApplicationDeploy2)
		v1.POST("/rds/:rds-name", rds.HandleRdsRequest)
		v1.GET("/sse", middleware.SseHeadersMiddleware(), func(ctx *gin.Context) {
			messages := []string{"a", "b", "c"}
			for _, msg := range messages {
				time.Sleep(time.Second * 1)
				ctx.SSEvent("message", msg)
				ctx.Writer.Flush()
			}
		})
		v1.POST("/deployAsync", middleware.SseHeadersMiddleware(), func(ctx *gin.Context) {
			deployer := deploy.NewDeployer(ctx)
			deployer.ApplicationDeployAsync()
		})
	}
}
