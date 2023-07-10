package main

import (
	"log"
    "github.com/gin-gonic/gin"

	api		"github.com/sw-maestro-kumofactory/miz-ball/api"
)



func main() {
	app := SetupRouter()
	log.Fatal(app.Run(":8080"))
	app.Run(":8080")
}

func SetupRouter() *gin.Engine {
	app := gin.Default()

	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	
	api.ApplyRoutes(app)
	return app
}