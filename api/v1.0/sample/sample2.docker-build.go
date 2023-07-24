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

	tarPath := "/app/repository/i-02b5064a1e36be086/repo.tar.gz"
	tags := []string{"434126037102.dkr.ecr.ap-northeast-2.amazonaws.com/kumo-customer:i-02b5064a1e36be086"}
	dockerfilePath := "coding-convention-sample-flask-f361fb2/Dockerfile"

	err = dockerclient.BuildImage(cli, tarPath, tags, dockerfilePath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func SAMPLE_TEST_BUILD2(c *gin.Context) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	tags := []string{"434126037102.dkr.ecr.ap-northeast-2.amazonaws.com/kumo-customer:i-02b5064a1e36be086"}
	bctx := "/app/repository/i-02b5064a1e36be086/coding-convention-sample-flask-f361fb2"

	err = dockerclient.BuildImage2(cli, bctx, tags)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
