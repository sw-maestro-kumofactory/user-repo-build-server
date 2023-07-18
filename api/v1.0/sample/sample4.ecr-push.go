package sample

import (
	"fmt"

	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"

	"github.com/sw-maestro-kumofactory/miz-ball/utils/ecr"
)

func SAMPLE_TEST_PUSH(c *gin.Context) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = ecr.Push(cli, "434126037102.dkr.ecr.ap-northeast-2.amazonaws.com/kumo-customer:i-02b5064a1e36be086")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
