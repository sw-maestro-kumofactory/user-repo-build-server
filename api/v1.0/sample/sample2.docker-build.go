package sample

import (
	"fmt"

	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"

	"github.com/sw-maestro-kumofactory/miz-ball/utils/dockerclient"
)

func SAMPLE_TEST_BUILD(c *gin.Context) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	tarPath := "sample-package/node-hello.tar.gz"
	tags := []string{"434126037102.dkr.ecr.ap-northeast-2.amazonaws.com/kumo-customer:i-02b5064a1e36be086"}
	dockerfilePath := "node-hello/Dockerfile"

	err = dockerclient.BuildImage(cli, tarPath, tags, dockerfilePath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
