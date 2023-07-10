package sample

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sw-maestro-kumofactory/miz-ball/utils/dockerfilegenerator"
	dfenum "github.com/sw-maestro-kumofactory/miz-ball/utils/dockerfilegenerator/enums"
)

func SAMPLE_TEST_CREATE(c *gin.Context) {
	builder := dockerfilegenerator.NewBuilder()

	builder.AddDirective(dfenum.FROM, "ubuntu:latest")
	builder.AddDirective(dfenum.RUN, "apt-get update && apt-get install -y curl")
	builder.AddDirective(dfenum.CMD, "echo 'Hello, Docker!'")

	err := builder.CreateDockerfile("./", "Dockerfile")
	if err != nil {
		fmt.Println("Failed to create Dockerfile:", err)
		return
	}

	fmt.Println("Dockerfile created successfully!")
}
